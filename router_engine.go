package krakend

import (
	"io"

	botdetector "github.com/devopsfaith/krakend-botdetector/gin"
	cors "github.com/devopsfaith/krakend-cors/gin"
	httpsecure "github.com/devopsfaith/krakend-httpsecure/gin"
	lua "github.com/devopsfaith/krakend-lua/router/gin"
	"github.com/devopsfaith/krakend/config"
	"github.com/devopsfaith/krakend/logging"
	"github.com/gin-gonic/gin"
)

// NewEngine creates a new gin engine with some default values and a secure middleware
func NewEngine(cfg config.ServiceConfig, logger logging.Logger, w io.Writer) *gin.Engine {
	if !cfg.Debug {
		gin.SetMode(gin.ReleaseMode)
	}

	engine := gin.New()
	engine.Use(gin.LoggerWithConfig(gin.LoggerConfig{Output: w}), gin.Recovery())

	engine.RedirectTrailingSlash = true
	engine.RedirectFixedPath = true
	engine.HandleMethodNotAllowed = true

	if corsMw := cors.New(cfg.ExtraConfig); corsMw != nil {
		engine.Use(corsMw)
	}

	if err := httpsecure.Register(cfg.ExtraConfig, engine); err != nil {
		logger.Warning(err)
	}

	lua.Register(logger, cfg.ExtraConfig, engine)

	botdetector.Register(cfg, logger, engine)

	return engine
}
