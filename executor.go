package krakend

import (
	"context"
	"fmt"
	"io"
	"os"

	krakendbf "github.com/devopsfaith/bloomfilter/krakend"
	cel "github.com/devopsfaith/krakend-cel"
	"github.com/devopsfaith/krakend-cobra"
	gelf "github.com/devopsfaith/krakend-gelf"
	"github.com/devopsfaith/krakend-gologging"
	"github.com/devopsfaith/krakend-jose"
	logstash "github.com/devopsfaith/krakend-logstash"
	metrics "github.com/devopsfaith/krakend-metrics/gin"
	opencensus "github.com/devopsfaith/krakend-opencensus"
	_ "github.com/devopsfaith/krakend-opencensus/exporter/influxdb"
	_ "github.com/devopsfaith/krakend-opencensus/exporter/jaeger"
	_ "github.com/devopsfaith/krakend-opencensus/exporter/prometheus"
	_ "github.com/devopsfaith/krakend-opencensus/exporter/stackdriver"
	_ "github.com/devopsfaith/krakend-opencensus/exporter/xray"
	_ "github.com/devopsfaith/krakend-opencensus/exporter/zipkin"
	_ "github.com/devopsfaith/krakend-opencensus/exporter/ocagent"
	pubsub "github.com/devopsfaith/krakend-pubsub"
	"github.com/devopsfaith/krakend-usage/client"
	"github.com/devopsfaith/krakend/config"
	"github.com/devopsfaith/krakend/logging"
	"github.com/devopsfaith/krakend/proxy"
	krakendrouter "github.com/devopsfaith/krakend/router"
	router "github.com/devopsfaith/krakend/router/gin"
	server "github.com/devopsfaith/krakend/transport/http/server/plugin"
	"github.com/gin-gonic/gin"
	"github.com/go-contrib/uuid"
	"github.com/letgoapp/krakend-influx"
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
// The service reguster function will register the service by the given name and port to all the available
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

// EngineFactory returns a gin engine, ready to be passed to the KrakenD RouterFactory
type EngineFactory interface {
	NewEngine(config.ServiceConfig, logging.Logger, io.Writer) *gin.Engine
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

	Middlewares []gin.HandlerFunc
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
			RunServer:      router.RunServerFunc(server.New(logger, krakendrouter.RunServer)),
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

	if err := influxdb.New(ctx, cfg.ExtraConfig, metricCollector, l); err != nil {
		l.Warning(err.Error())
	}

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
