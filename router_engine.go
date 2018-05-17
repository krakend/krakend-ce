package main

import (
	"time"

	httpsecure "github.com/devopsfaith/krakend-httpsecure/gin"
	metricsgin "github.com/devopsfaith/krakend-metrics/gin"
	"github.com/devopsfaith/krakend/config"
	"github.com/devopsfaith/krakend/logging"
	"github.com/gin-gonic/gin"
)

// NewEngine creates a new gin engine with some default values, a secure middleware and an stats endpoint
func NewEngine(cfg config.ServiceConfig, logger logging.Logger, metricCollector *metricsgin.Metrics) *gin.Engine {
	engine := gin.New()

	engine.RedirectTrailingSlash = true
	engine.RedirectFixedPath = true
	engine.HandleMethodNotAllowed = true

	engine.Use(gin.Recovery(), GinLogger(logger))

	if err := httpsecure.Register(cfg.ExtraConfig, engine); err != nil {
		logger.Error(err)
	}

	engine.GET("/__stats/", metricCollector.NewExpHandler())

	return engine
}

func GinLogger(logger logging.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		c.Next()

		end := time.Now()

		if raw != "" {
			path = path + "?" + raw
		}

		event := map[string]interface{}{
			"path":      path,
			"method":    c.Request.Method,
			"latency":   end.Sub(start).String(),
			"status":    c.Writer.Status(),
			"client-ip": c.ClientIP(),
		}

		if c.Errors != nil {
			event["error"] = c.Errors.Errors()
		}

		logger.Debug("http request served", event)
	}
}
