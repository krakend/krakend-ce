package main

import (
	httpsecure "github.com/devopsfaith/krakend-httpsecure/gin"
	metricsgin "github.com/devopsfaith/krakend-metrics/gin"
	"github.com/devopsfaith/krakend/config"
	"github.com/devopsfaith/krakend/logging"
	"github.com/gin-gonic/gin"
)

// NewEngine creates a new gin engine with some default values, a secure middleware and an stats endpoint
func NewEngine(cfg config.ServiceConfig, logger logging.Logger, metricCollector *metricsgin.Metrics) *gin.Engine {
	engine := gin.Default()
	engine.RedirectTrailingSlash = true
	engine.RedirectFixedPath = true
	engine.HandleMethodNotAllowed = true
	if err := httpsecure.Register(cfg.ExtraConfig, engine); err != nil {
		logger.Error(err)
	}
	engine.GET("/__stats/", metricCollector.NewExpHandler())
	return engine
}
