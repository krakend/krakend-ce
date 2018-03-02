package main

import (
	"context"
	"os"
	"time"

	"github.com/devopsfaith/krakend-cobra"
	"github.com/devopsfaith/krakend-gologging"
	metrics "github.com/devopsfaith/krakend-metrics/gin"
	"github.com/devopsfaith/krakend/config"
	"github.com/devopsfaith/krakend/logging"
	router "github.com/devopsfaith/krakend/router/gin"
	"github.com/gin-gonic/gin"
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

		RegisterSubscriberFactories(ctx, cfg, logger)

		// create the metrics collector
		metricCollector := metrics.New(ctx, time.Minute, logger)

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
