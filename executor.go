package krakend

import (
	"context"
	"fmt"
	krakendbf "github.com/devopsfaith/bloomfilter/krakend"
	cel "github.com/devopsfaith/krakend-cel"
	cmd "github.com/devopsfaith/krakend-cobra"
	cors "github.com/devopsfaith/krakend-cors/mux"
	gelf "github.com/devopsfaith/krakend-gelf"
	gologging "github.com/devopsfaith/krakend-gologging"
	httpsecure "github.com/devopsfaith/krakend-httpsecure/mux"
	jose "github.com/devopsfaith/krakend-jose"
	logstash "github.com/devopsfaith/krakend-logstash"
	lua "github.com/devopsfaith/krakend-lua/router/mux"
	metrics "github.com/devopsfaith/krakend-metrics/mux"
	opencensus "github.com/devopsfaith/krakend-opencensus"
	_ "github.com/devopsfaith/krakend-opencensus/exporter/datadog"
	_ "github.com/devopsfaith/krakend-opencensus/exporter/influxdb"
	_ "github.com/devopsfaith/krakend-opencensus/exporter/jaeger"
	_ "github.com/devopsfaith/krakend-opencensus/exporter/ocagent"
	_ "github.com/devopsfaith/krakend-opencensus/exporter/prometheus"
	_ "github.com/devopsfaith/krakend-opencensus/exporter/stackdriver"
	_ "github.com/devopsfaith/krakend-opencensus/exporter/xray"
	_ "github.com/devopsfaith/krakend-opencensus/exporter/zipkin"
	pubsub "github.com/devopsfaith/krakend-pubsub"
	"github.com/devopsfaith/krakend-usage/client"
	"github.com/devopsfaith/krakend/config"
	"github.com/devopsfaith/krakend/logging"
	"github.com/devopsfaith/krakend/proxy"
	krakendrouter "github.com/devopsfaith/krakend/router"
	"github.com/devopsfaith/krakend/router/httptreemux"
	router "github.com/devopsfaith/krakend/router/mux"
	server "github.com/devopsfaith/krakend/transport/http/server/plugin"
	"github.com/go-contrib/uuid"
	"io"
	"net/http"
	"os"
)

// NewExecutor returns an executor for the cmd package. The executor initalizes the entire gateway by
// registering the components and composing a RouterFactory wrapping all the middlewares.
func NewExecutor(ctx context.Context) cmd.Executor {
	eb := new(ExecutorBuilder)
	return eb.NewCmdExecutor(ctx)
}

// PluginLoader defines the interface for the collaborator responsible of starting the plugin loaders
type PluginLoader interface {
	Load(folder, pattern string, logger logging.Logger)
}

// SubscriberFactoriesRegister registers all the required subscriber factories from the available service
// discover components and adapters and returns a service register function.
// The service register function will register the service by the given name and port to all the available
// service discover clients
type SubscriberFactoriesRegister interface {
	Register(context.Context, config.ServiceConfig, logging.Logger) func(string, int)
}

// TokenRejecterFactory returns a jose.ChainedRejecterFactory containing all the required jose.RejecterFactory.
// It also should setup and manage any service related to the management of the revocation process, if required.
type TokenRejecterFactory interface {
	NewTokenRejecter(context.Context, config.ServiceConfig, logging.Logger, func(string, int)) (jose.ChainedRejecterFactory, error)
}

// MetricsAndTracesRegister registers the defined observability components and returns a metrics collector,
// if required.
type MetricsAndTracesRegister interface {
	Register(context.Context, config.ServiceConfig, logging.Logger) *metrics.Metrics
}

// EngineFactory returns a mux engine, ready to be passed to the KrakenD RouterFactory
type EngineFactory interface {
	NewEngine(config.ServiceConfig, logging.Logger, io.Writer) httptreemux.Engine
}

// ProxyFactory returns a KrakenD proxy factory, ready to be passed to the KrakenD RouterFactory
type ProxyFactory interface {
	NewProxyFactory(logging.Logger, proxy.BackendFactory, *metrics.Metrics) proxy.Factory
}

// BackendFactory returns a KrakenD backend factory, ready to be passed to the KrakenD proxy factory
type BackendFactory interface {
	NewBackendFactory(context.Context, logging.Logger, *metrics.Metrics) proxy.BackendFactory
}

// HandlerFactory returns a KrakenD router handler factory, ready to be passed to the KrakenD RouterFactory
type HandlerFactory interface {
	NewHandlerFactory(logging.Logger, *metrics.Metrics, jose.RejecterFactory) router.HandlerFactory
}

// LoggerFactory returns a KrakenD Logger factory, ready to be passed to the KrakenD RouterFactory
type LoggerFactory interface {
	NewLogger(config.ServiceConfig) (logging.Logger, io.Writer, error)
}

// RunServer defines the interface of a function used by the KrakenD router to start the service
type RunServer func(context.Context, config.ServiceConfig, http.Handler) error

// RunServerFactory returns a RunServer with several wraps around the injected one
type RunServerFactory interface {
	NewRunServer(logging.Logger, router.RunServerFunc) RunServer
}

// ExecutorBuilder is a composable builder. Every injected property is used by the NewCmdExecutor method.
type ExecutorBuilder struct {
	LoggerFactory               LoggerFactory
	PluginLoader                PluginLoader
	SubscriberFactoriesRegister SubscriberFactoriesRegister
	TokenRejecterFactory        TokenRejecterFactory
	MetricsAndTracesRegister    MetricsAndTracesRegister
	EngineFactory               EngineFactory
	ProxyFactory                ProxyFactory
	BackendFactory              BackendFactory
	HandlerFactory              HandlerFactory
	RunServerFactory            RunServerFactory

	Middlewares []router.HandlerMiddleware
}

// NewCmdExecutor returns an executor for the cmd package. The executor initalizes the entire gateway by
// delegating most of the tasks to the injected collaborators. They register the components and
// compose a RouterFactory wrapping all the middlewares.
// Every nil collaborator is replaced by the default one offered by this package.
func (e *ExecutorBuilder) NewCmdExecutor(ctx context.Context) cmd.Executor {
	e.checkCollaborators()

	return func(cfg config.ServiceConfig) {
		logger, gelfWriter, gelfErr := e.LoggerFactory.NewLogger(cfg)
		if gelfErr != nil {
			return
		}

		logger.Info("Listening on port:", cfg.Port)

		startReporter(ctx, logger, cfg)

		if cfg.Plugin != nil {
			e.PluginLoader.Load(cfg.Plugin.Folder, cfg.Plugin.Pattern, logger)
		}

		metricCollector := e.MetricsAndTracesRegister.Register(ctx, cfg, logger)

		tokenRejecterFactory, err := e.TokenRejecterFactory.NewTokenRejecter(
			ctx,
			cfg,
			logger,
			e.SubscriberFactoriesRegister.Register(ctx, cfg, logger),
		)
		if err != nil {
			logger.Warning("bloomFilter:", err.Error())
		}

		e.Middlewares = lua.RegisterMiddleware(logger, cfg.ExtraConfig, httptreemux.ParamsExtractor, e.Middlewares)
		e.Middlewares = append(e.Middlewares, httpsecure.NewSecureMw(cfg.ExtraConfig))
		e.Middlewares = Register(cfg.ExtraConfig, logger, e.Middlewares)

		// setup the krakend router
		routerFactory := router.NewFactory(router.Config{
			Engine: e.EngineFactory.NewEngine(cfg, logger, gelfWriter),
			ProxyFactory: e.ProxyFactory.NewProxyFactory(
				logger,
				e.BackendFactory.NewBackendFactory(ctx, logger, metricCollector),
				metricCollector,
			),
			Middlewares:    e.Middlewares,
			Logger:         logger,
			HandlerFactory: e.HandlerFactory.NewHandlerFactory(logger, metricCollector, tokenRejecterFactory),
			RunServer:      router.RunServerFunc(e.RunServerFactory.NewRunServer(logger, krakendrouter.RunServer)),
		})

		// start the engines
		routerFactory.NewWithContext(ctx).Run(cfg)
	}
}

func (e *ExecutorBuilder) checkCollaborators() {
	if e.PluginLoader == nil {
		e.PluginLoader = new(pluginLoader)
	}
	if e.SubscriberFactoriesRegister == nil {
		e.SubscriberFactoriesRegister = new(registerSubscriberFactories)
	}
	if e.TokenRejecterFactory == nil {
		e.TokenRejecterFactory = new(BloomFilterJWT)
	}
	if e.MetricsAndTracesRegister == nil {
		e.MetricsAndTracesRegister = new(MetricsAndTraces)
	}
	if e.EngineFactory == nil {
		e.EngineFactory = new(engineFactory)
	}
	if e.ProxyFactory == nil {
		e.ProxyFactory = new(proxyFactory)
	}
	if e.BackendFactory == nil {
		e.BackendFactory = new(backendFactory)
	}
	if e.HandlerFactory == nil {
		e.HandlerFactory = new(handlerFactory)
	}
	if e.LoggerFactory == nil {
		e.LoggerFactory = new(LoggerBuilder)
	}
	if e.RunServerFactory == nil {
		e.RunServerFactory = new(DefaultRunServerFactory)
	}
}

// DefaultRunServerFactory creates the default RunServer by wrapping the injected RunServer
// with the plugin loader and the CORS module
type DefaultRunServerFactory struct{}

func (d *DefaultRunServerFactory) NewRunServer(l logging.Logger, next router.RunServerFunc) RunServer {
	return RunServer(server.New(
		l,
		server.RunServer(NewRunServer(NewRunServerWithLogger(RunServer(next), l))),
	))
}

// NewRunServer returns a RunServer wrapping the injected one with a CORS middleware, so it is called before the
// actual router checks the URL, method and other details related to selecting the proper handler for the
// incoming request
func NewRunServer(next RunServer) RunServer {
	return NewRunServerWithLogger(next, nil)
}

// NewRunServerWithLogger returns a RunServer wrapping the injected one with a CORS middleware, so it is called before the
// actual router checks the URL, method and other details related to selecting the proper handler for the
// incoming request
func NewRunServerWithLogger(next RunServer, l logging.Logger) RunServer {
	return func(ctx context.Context, cfg config.ServiceConfig, handler http.Handler) error {
		corsMw := cors.NewWithLogger(cfg.ExtraConfig, l)
		if corsMw == nil {
			return next(ctx, cfg, handler)
		}
		return next(ctx, cfg, corsMw.Handler(handler))
	}
}

// LoggerBuilder is the default BuilderFactory implementation.
type LoggerBuilder struct{}

// NewLogger sets up the logging components as defined at the configuration.
func (LoggerBuilder) NewLogger(cfg config.ServiceConfig) (logging.Logger, io.Writer, error) {
	var writers []io.Writer
	gelfWriter, gelfErr := gelf.NewWriter(cfg.ExtraConfig)
	if gelfErr == nil {
		writers = append(writers, gelfWriterWrapper{gelfWriter})
		gologging.SetFormatterSelector(func(w io.Writer) string {
			switch w.(type) {
			case gelfWriterWrapper:
				return "%{message}"
			default:
				return gologging.DefaultPattern
			}
		})
	}
	logger, gologgingErr := logstash.NewLogger(cfg.ExtraConfig)

	if gologgingErr != nil {
		logger, gologgingErr = gologging.NewLogger(cfg.ExtraConfig, writers...)

		if gologgingErr != nil {
			var err error
			logger, err = logging.NewLogger("DEBUG", os.Stdout, "")
			if err != nil {
				return logger, gelfWriter, err
			}
			logger.Error("unable to create the gologging logger:", gologgingErr.Error())
		}
	}
	if gelfErr != nil {
		logger.Error("unable to create the GELF writer:", gelfErr.Error())
	}
	return logger, gelfWriter, nil
}

// BloomFilterJWT is the default TokenRejecterFactory implementation.
type BloomFilterJWT struct{}

// NewTokenRejecter registers the bloomfilter component and links it to a token rejecter. Then it returns a chained
// rejecter factory with the created token rejecter and other based on the CEL component.
func (t BloomFilterJWT) NewTokenRejecter(ctx context.Context, cfg config.ServiceConfig, l logging.Logger, reg func(n string, p int)) (jose.ChainedRejecterFactory, error) {
	rejecter, err := krakendbf.Register(ctx, "krakend-bf", cfg, l, reg)

	return jose.ChainedRejecterFactory([]jose.RejecterFactory{
		jose.RejecterFactoryFunc(func(_ logging.Logger, _ *config.EndpointConfig) jose.Rejecter {
			return jose.RejecterFunc(rejecter.RejectToken)
		}),
		jose.RejecterFactoryFunc(func(l logging.Logger, cfg *config.EndpointConfig) jose.Rejecter {
			if r := cel.NewRejecter(l, cfg); r != nil {
				return r
			}
			return jose.FixedRejecter(false)
		}),
	}), err
}

// MetricsAndTraces is the default implementation of the MetricsAndTracesRegister interface.
type MetricsAndTraces struct{}

// Register registers the metrcis, influx and opencensus packages as required by the given configuration.
func (MetricsAndTraces) Register(ctx context.Context, cfg config.ServiceConfig, l logging.Logger) *metrics.Metrics {
	metricCollector := metrics.New(ctx, cfg.ExtraConfig, l)

	if err := opencensus.Register(ctx, cfg, append(opencensus.DefaultViews, pubsub.OpenCensusViews...)...); err != nil {
		l.Warning("opencensus:", err.Error())
	}

	return metricCollector
}

const (
	usageDisable = "USAGE_DISABLE"
)

func startReporter(ctx context.Context, logger logging.Logger, cfg config.ServiceConfig) {
	if os.Getenv(usageDisable) == "1" {
		logger.Info("usage report client disabled")
		return
	}

	clusterID, err := cfg.Hash()
	if err != nil {
		logger.Warning("unable to hash the service configuration:", err.Error())
		return
	}

	go func() {
		serverID := uuid.NewV4().String()
		logger.Info(fmt.Sprintf("registering usage stats for cluster ID '%s'", clusterID))

		if err := client.StartReporter(ctx, client.Options{
			ClusterID: clusterID,
			ServerID:  serverID,
		}); err != nil {
			logger.Warning("unable to create the usage report client:", err.Error())
		}
	}()
}

type gelfWriterWrapper struct {
	io.Writer
}
