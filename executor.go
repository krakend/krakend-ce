package main

import (
	"context"
	"fmt"
	"os"

	krakendbf "github.com/devopsfaith/bloomfilter/krakend"
	"github.com/devopsfaith/krakend-cobra"
	"github.com/devopsfaith/krakend-gologging"
	jose "github.com/devopsfaith/krakend-jose"
	metrics "github.com/devopsfaith/krakend-metrics/gin"
	opencensus "github.com/devopsfaith/krakend-opencensus"
	_ "github.com/devopsfaith/krakend-opencensus/exporter/influxdb"
	_ "github.com/devopsfaith/krakend-opencensus/exporter/jaeger"
	_ "github.com/devopsfaith/krakend-opencensus/exporter/prometheus"
	_ "github.com/devopsfaith/krakend-opencensus/exporter/zipkin"
	"github.com/devopsfaith/krakend-usage/client"
	"github.com/devopsfaith/krakend/config"
	"github.com/devopsfaith/krakend/logging"
	krakendrouter "github.com/devopsfaith/krakend/router"
	router "github.com/devopsfaith/krakend/router/gin"
	"github.com/gin-gonic/gin"
	"github.com/go-contrib/uuid"
	"github.com/letgoapp/krakend-influx"
)

func NewExecutor(ctx context.Context) cmd.Executor {
	return func(cfg config.ServiceConfig) {
		logger, gologgingErr := gologging.NewLogger(cfg.ExtraConfig)
		if gologgingErr != nil {
			var err error
			logger, err = logging.NewLogger("DEBUG", os.Stdout, "")
			if err != nil {
				return
			}
			logger.Error("unable to create the gologgin logger:", gologgingErr.Error())
		}

		startReporter(ctx, logger, cfg)

		reg := RegisterSubscriberFactories(ctx, cfg, logger)

		// create the metrics collector
		metricCollector := metrics.New(ctx, cfg.ExtraConfig, logger)

		if err := influxdb.New(ctx, cfg.ExtraConfig, metricCollector, logger); err != nil {
			logger.Error(err.Error())
		}

		if err := opencensus.Register(ctx, cfg); err != nil {
			logger.Error("opencensus:", err.Error())
		}

		rejecter, err := krakendbf.Register(ctx, "krakend-bf", cfg, logger, reg)
		// rejecter, err := krakendbf.Register(ctx, cfg.Name, cfg, logger, reg)
		if err != nil {
			logger.Error("registering the BloomFilter:", err.Error())
		}

		// setup the krakend router
		routerFactory := router.NewFactory(router.Config{
			Engine:         NewEngine(cfg, logger),
			ProxyFactory:   NewProxyFactory(logger, NewBackendFactory(logger, metricCollector), metricCollector),
			Middlewares:    []gin.HandlerFunc{},
			Logger:         logger,
			HandlerFactory: NewHandlerFactory(logger, metricCollector, jose.RejecterFunc(rejecter.RejectToken)),
			RunServer:      krakendrouter.RunServer,
		})

		// start the engines
		routerFactory.NewWithContext(ctx).Run(cfg)
	}
}

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
