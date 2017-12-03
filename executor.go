package main

import (
	"context"
	"os"
	"time"

	"github.com/devopsfaith/krakend-cobra"
	gologging "github.com/devopsfaith/krakend-gologging"
	metrics "github.com/devopsfaith/krakend-metrics/gin"
	"github.com/devopsfaith/krakend/config"
	"github.com/devopsfaith/krakend/logging"
	router "github.com/devopsfaith/krakend/router/gin"
	"github.com/gin-gonic/gin"
)

func NewExecutor(ctx context.Context) cmd.Executor {
	return func(cfg config.ServiceConfig) {
		logger, err := gologging.NewLogger(cfg.ExtraConfig)
		if err != nil {
			logger, err = logging.NewLogger("DEBUG", os.Stdout, "")
			if err != nil {
				return
			}
			logger.Error("unable to create the gologgin logger:", err.Error())
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
			HandlerFactory: NewHandlerFactory(metricCollector),
		})

		// start the engines
		routerFactory.NewWithContext(ctx).Run(cfg)
	}
}
