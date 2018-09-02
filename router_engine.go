package main

import (
	cors "github.com/devopsfaith/krakend-cors/gin"
	httpsecure "github.com/devopsfaith/krakend-httpsecure/gin"
	"github.com/devopsfaith/krakend/config"
	"github.com/devopsfaith/krakend/logging"
	"github.com/gin-gonic/gin"
)

// NewEngine creates a new gin engine with some default values and a secure middleware
func NewEngine(cfg config.ServiceConfig, logger logging.Logger) *gin.Engine {
	engine := gin.Default()

	engine.RedirectTrailingSlash = true
	engine.RedirectFixedPath = true
	engine.HandleMethodNotAllowed = true

	if mw := cors.New(cfg.ExtraConfig); mw != nil {
		engine.Use(mw)
	}

	if err := httpsecure.Register(cfg.ExtraConfig, engine); err != nil {
		logger.Error(err)
	}

	return engine
}
