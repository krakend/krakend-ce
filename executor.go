package krakend

import (
	"context"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-contrib/uuid"
	"golang.org/x/sync/errgroup"

	asyncamqp "github.com/devopsfaith/krakend-amqp/v2/async"
	gelf "github.com/devopsfaith/krakend-gelf/v2"
	influxdb "github.com/devopsfaith/krakend-influx/v2"
	metrics "github.com/devopsfaith/krakend-metrics/v2/gin"
	krakendbf "github.com/krakendio/bloomfilter/v2/krakend"
	cel "github.com/krakendio/krakend-cel/v2"
	cmd "github.com/krakendio/krakend-cobra/v2"
	cors "github.com/krakendio/krakend-cors/v2/gin"
	gologging "github.com/krakendio/krakend-gologging/v2"
	jose "github.com/krakendio/krakend-jose/v2"
	logstash "github.com/krakendio/krakend-logstash/v2"
	opencensus "github.com/krakendio/krakend-opencensus/v2"
	_ "github.com/krakendio/krakend-opencensus/v2/exporter/datadog"
	_ "github.com/krakendio/krakend-opencensus/v2/exporter/influxdb"
	_ "github.com/krakendio/krakend-opencensus/v2/exporter/jaeger"
	_ "github.com/krakendio/krakend-opencensus/v2/exporter/ocagent"
	_ "github.com/krakendio/krakend-opencensus/v2/exporter/prometheus"
	_ "github.com/krakendio/krakend-opencensus/v2/exporter/stackdriver"
	_ "github.com/krakendio/krakend-opencensus/v2/exporter/xray"
	_ "github.com/krakendio/krakend-opencensus/v2/exporter/zipkin"
	pubsub "github.com/krakendio/krakend-pubsub/v2"
	"github.com/krakendio/krakend-usage/client"
	"github.com/luraproject/lura/v2/async"
	"github.com/luraproject/lura/v2/config"
	"github.com/luraproject/lura/v2/core"
	"github.com/luraproject/lura/v2/logging"
	"github.com/luraproject/lura/v2/proxy"
	router "github.com/luraproject/lura/v2/router/gin"
	serverhttp "github.com/luraproject/lura/v2/transport/http/server"
	server "github.com/luraproject/lura/v2/transport/http/server/plugin"
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

// EngineFactory returns a gin engine, ready to be passed to the KrakenD RouterFactory
type EngineFactory interface {
	NewEngine(config.ServiceConfig, router.EngineOptions) *gin.Engine
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

// AgentStarter defines a type that starts a set of agents
type AgentStarter interface {
	Start(
		context.Context,
		[]*config.AsyncAgent,
		logging.Logger,
		chan<- string,
		proxy.Factory,
	) func() error
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
	AgentStarterFactory         AgentStarter

	Middlewares []gin.HandlerFunc
}

// NewCmdExecutor returns an executor for the cmd package. The executor initializes the entire gateway by
// delegating most of the tasks to the injected collaborators. They register the components and
// compose a RouterFactory wrapping all the middlewares.
// Every nil collaborator is replaced by the default one offered by this package.
func (e *ExecutorBuilder) NewCmdExecutor(ctx context.Context) cmd.Executor {
	e.checkCollaborators()

	return func(cfg config.ServiceConfig) {
		cfg.Normalize()

		logger, gelfWriter, gelfErr := e.LoggerFactory.NewLogger(cfg)
		if gelfErr != nil {
			return
		}

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
		if err != nil && err != krakendbf.ErrNoConfig {
			logger.Warning("[SERVICE: Bloomfilter]", err.Error())
		}

		pf := e.ProxyFactory.NewProxyFactory(
			logger,
			e.BackendFactory.NewBackendFactory(ctx, logger, metricCollector),
			metricCollector,
		)

		agentPing := make(chan string, len(cfg.AsyncAgents))

		// setup the krakend router
		routerFactory := router.NewFactory(router.Config{
			Engine: e.EngineFactory.NewEngine(cfg, router.EngineOptions{
				Logger: logger,
				Writer: gelfWriter,
				Health: (<-chan string)(agentPing),
			}),
			ProxyFactory:   pf,
			Middlewares:    e.Middlewares,
			Logger:         logger,
			HandlerFactory: e.HandlerFactory.NewHandlerFactory(logger, metricCollector, tokenRejecterFactory),
			RunServer:      router.RunServerFunc(e.RunServerFactory.NewRunServer(logger, serverhttp.RunServer)),
		})

		// start the engines
		logger.Info("Starting the KrakenD instance")

		if len(cfg.AsyncAgents) == 0 {
			routerFactory.NewWithContext(ctx).Run(cfg)
			return
		}

		// start the async agents in the same error group as the router
		g, gctx := errgroup.WithContext(ctx)
		gctx, closeGroupCtx := context.WithCancel(gctx)

		if cfg.SequentialStart {
			waitAgents := e.AgentStarterFactory.Start(gctx, cfg.AsyncAgents, logger, (chan<- string)(agentPing), pf)
			g.Go(waitAgents)
		} else {
			g.Go(func() error {
				return e.AgentStarterFactory.Start(gctx, cfg.AsyncAgents, logger, (chan<- string)(agentPing), pf)()
			})
		}

		g.Go(func() error {
			logger.Info("[SERVICE: Gin] Building the router")
			routerFactory.NewWithContext(ctx).Run(cfg)
			closeGroupCtx()
			return nil
		})

		g.Wait()
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
	if e.AgentStarterFactory == nil {
		e.AgentStarterFactory = async.AgentStarter([]async.Factory{asyncamqp.StartAgent})
	}
}

// DefaultRunServerFactory creates the default RunServer by wrapping the injected RunServer
// with the plugin loader and the CORS module
type DefaultRunServerFactory struct{}

func (*DefaultRunServerFactory) NewRunServer(l logging.Logger, next router.RunServerFunc) RunServer {
	return RunServer(server.New(
		l,
		server.RunServer(cors.NewRunServer(cors.NewRunServerWithLogger(cors.RunServer(next), l))),
	))
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
			logger, err = logging.NewLogger("DEBUG", os.Stdout, "KRAKEND")
			if err != nil {
				return logger, gelfWriter, err
			}
			logger.Error("[SERVICE: Logging] Unable to create the logger:", gologgingErr.Error())
		}
	}
	if gelfErr != nil && gelfErr != gelf.ErrWrongConfig {
		logger.Error("[SERVICE: Logging][GELF] Unable to create the writer:", gelfErr.Error())
	}
	return logger, gelfWriter, nil
}

// BloomFilterJWT is the default TokenRejecterFactory implementation.
type BloomFilterJWT struct{}

// NewTokenRejecter registers the bloomfilter component and links it to a token rejecter. Then it returns a chained
// rejecter factory with the created token rejecter and other based on the CEL component.
func (BloomFilterJWT) NewTokenRejecter(ctx context.Context, cfg config.ServiceConfig, l logging.Logger, reg func(n string, p int)) (jose.ChainedRejecterFactory, error) {
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

// Register registers the metrics, influx and opencensus packages as required by the given configuration.
func (MetricsAndTraces) Register(ctx context.Context, cfg config.ServiceConfig, l logging.Logger) *metrics.Metrics {
	metricCollector := metrics.New(ctx, cfg.ExtraConfig, l)

	if err := influxdb.New(ctx, cfg.ExtraConfig, metricCollector, l); err != nil {
		if err != influxdb.ErrNoConfig {
			l.Warning("[SERVICE: InfluxDB]", err.Error())
		}
	} else {
		l.Debug("[SERVICE: InfluxDB] Service correctly registered")
	}

	if err := opencensus.Register(ctx, cfg, append(opencensus.DefaultViews, pubsub.OpenCensusViews...)...); err != nil {
		if err != opencensus.ErrNoConfig {
			l.Warning("[SERVICE: OpenCensus]", err.Error())
		}
	} else {
		l.Debug("[SERVICE: OpenCensus] Service correctly registered")
	}

	return metricCollector
}

const (
	usageDisable = "USAGE_DISABLE"
	usageDelay   = 5 * time.Second
)

func startReporter(ctx context.Context, logger logging.Logger, cfg config.ServiceConfig) {
	logPrefix := "[SERVICE: Telemetry]"
	if os.Getenv(usageDisable) == "1" {
		return
	}

	clusterID, err := cfg.Hash()
	if err != nil {
		logger.Debug(logPrefix, "Unable to create the Cluster ID hash:", err.Error())
		return
	}

	go func() {
		time.Sleep(usageDelay)

		serverID := uuid.NewV4().String()
		logger.Debug(logPrefix, "Registering usage stats for Cluster ID", clusterID)

		if err := client.StartReporter(ctx, client.Options{
			ClusterID: clusterID,
			ServerID:  serverID,
			Version:   core.KrakendVersion,
		}); err != nil {
			logger.Debug(logPrefix, "Unable to create the usage report client:", err.Error())
		}
	}()
}

type gelfWriterWrapper struct {
	io.Writer
}
