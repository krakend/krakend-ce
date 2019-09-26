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
	"github.com/devopsfaith/krakend-opencensus"
	_ "github.com/devopsfaith/krakend-opencensus/exporter/influxdb"
	_ "github.com/devopsfaith/krakend-opencensus/exporter/jaeger"
	_ "github.com/devopsfaith/krakend-opencensus/exporter/prometheus"
	_ "github.com/devopsfaith/krakend-opencensus/exporter/stackdriver"
	_ "github.com/devopsfaith/krakend-opencensus/exporter/xray"
	_ "github.com/devopsfaith/krakend-opencensus/exporter/zipkin"
	pubsub "github.com/devopsfaith/krakend-pubsub"
	"github.com/devopsfaith/krakend-usage/client"
	"github.com/devopsfaith/krakend/config"
	"github.com/devopsfaith/krakend/logging"
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
	return func(cfg config.ServiceConfig) {
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
					return
				}
				logger.Error("unable to create the gologging logger:", gologgingErr.Error())
			}
		}
		if gelfErr != nil {
			logger.Error("unable to create the GELF writer:", gelfErr.Error())
		}

		logger.Info("Listening on port:", cfg.Port)

		startReporter(ctx, logger, cfg)

		if cfg.Plugin != nil {
			LoadPlugins(cfg.Plugin.Folder, cfg.Plugin.Pattern, logger)
		}

		reg := RegisterSubscriberFactories(ctx, cfg, logger)

		// create the metrics collector
		metricCollector := metrics.New(ctx, cfg.ExtraConfig, logger)

		if err := influxdb.New(ctx, cfg.ExtraConfig, metricCollector, logger); err != nil {
			logger.Warning(err.Error())
		}

		if err := opencensus.Register(ctx, cfg, append(opencensus.DefaultViews, pubsub.OpenCensusViews...)...); err != nil {
			logger.Warning("opencensus:", err.Error())
		}

		rejecter, err := krakendbf.Register(ctx, "krakend-bf", cfg, logger, reg)
		if err != nil {
			logger.Warning("bloomFilter:", err.Error())
		}

		tokenRejecterFactory := jose.ChainedRejecterFactory([]jose.RejecterFactory{
			jose.RejecterFactoryFunc(func(_ logging.Logger, _ *config.EndpointConfig) jose.Rejecter {
				return jose.RejecterFunc(rejecter.RejectToken)
			}),
			jose.RejecterFactoryFunc(func(l logging.Logger, cfg *config.EndpointConfig) jose.Rejecter {
				if r := cel.NewRejecter(l, cfg); r != nil {
					return r
				}
				return jose.FixedRejecter(false)
			}),
		})

		// setup the krakend router
		routerFactory := router.NewFactory(router.Config{
			Engine: NewEngine(cfg, logger, gelfWriter),
			ProxyFactory: NewProxyFactory(
				logger, NewBackendFactoryWithContext(ctx, logger, metricCollector),
				metricCollector,
			),
			Middlewares:    []gin.HandlerFunc{},
			Logger:         logger,
			HandlerFactory: NewHandlerFactory(logger, metricCollector, tokenRejecterFactory),
			RunServer:      router.RunServerFunc(server.New(logger, krakendrouter.RunServer)),
		})

		// start the engines
		routerFactory.NewWithContext(ctx).Run(cfg)
	}
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
