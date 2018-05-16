package main

import (
	"context"
	"os"
	"time"

	"github.com/devopsfaith/krakend-cobra"
	"github.com/devopsfaith/krakend-logstash"
	metrics "github.com/devopsfaith/krakend-metrics/gin"
	"github.com/devopsfaith/krakend/config"
	"github.com/devopsfaith/krakend/logging"
	router "github.com/devopsfaith/krakend/router/gin"
	"github.com/gin-gonic/gin"
	"github.com/letgoapp/krakend-influx"
)

func NewExecutor(ctx context.Context) cmd.Executor {
	return func(cfg config.ServiceConfig) {
		var logger logging.Logger
		logger, gologgingErr := logstsash.NewLogger(cfg.ExtraConfig)
		if gologgingErr != nil {
			var err error
			logger, err = logging.NewLogger("DEBUG", os.Stdout, "")
			if err != nil {
				return
			}
			logger.Error("unable to create the gologgin logger:", gologgingErr.Error())
		}

		RegisterSubscriberFactories(ctx, cfg, logger)

		// create the metrics collector
		metricCollector := metrics.New(ctx, time.Minute, logger)

		if err := influxdb.New(ctx, cfg.ExtraConfig, metricCollector, logger); err != nil {
			logger.Error(err.Error())
			return
		}

		// setup the krakend router
		routerFactory := router.NewFactory(router.Config{
			Engine:         NewEngine(cfg, logger, metricCollector),
			ProxyFactory:   NewProxyFactory(logger, NewBackendFactory(logger, metricCollector), metricCollector),
			Middlewares:    []gin.HandlerFunc{},
			Logger:         logger,
			HandlerFactory: NewHandlerFactory(logger, metricCollector),
		})

		// start the engines
		routerFactory.NewWithContext(ctx).Run(cfg)
	}
}
