package main

import (
	"context"
	"os"
	"time"

	"github.com/devopsfaith/krakend-cobra"
	gologging "github.com/devopsfaith/krakend-gologging"
	"github.com/devopsfaith/krakend-martian"
	metrics "github.com/devopsfaith/krakend-metrics/gin"
	"github.com/devopsfaith/krakend/config"
	"github.com/devopsfaith/krakend/logging"
	"github.com/devopsfaith/krakend/plugin"
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

		if "" != os.Getenv("KRAKEND_ENABLE_PLUGINS") && cfg.Plugin != nil {
			logger.Info("Plugin experiment enabled!")
			pluginsLoaded, err := plugin.Load(*cfg.Plugin)
			if err != nil {
				logger.Error(err.Error())
			}
			logger.Info("Total plugins loaded:", pluginsLoaded)
		}

		RegisterSubscriberFactories(ctx, cfg, logger)

		// create the metrics collector
		metricCollector := metrics.New(ctx, time.Minute, logger)

		martian.Register()

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
