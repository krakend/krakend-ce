package krakend

import (
	"github.com/gin-gonic/gin"

	botdetector "github.com/devopsfaith/krakend-botdetector/v2/gin"
	httpsecure "github.com/devopsfaith/krakend-httpsecure/v2/gin"
	lua "github.com/devopsfaith/krakend-lua/v2/router/gin"
	opencensus "github.com/devopsfaith/krakend-opencensus/v2/router/gin"
	"github.com/luraproject/lura/v2/config"
	"github.com/luraproject/lura/v2/core"
	luragin "github.com/luraproject/lura/v2/router/gin"
	"github.com/luraproject/lura/v2/transport/http/server"
)

// NewEngine creates a new gin engine with some default values and a secure middleware
func NewEngine(cfg config.ServiceConfig, opt luragin.EngineOptions) *gin.Engine {
	engine := luragin.NewEngine(cfg, opt)

	engine.NoRoute(opencensus.HandlerFunc(&config.EndpointConfig{Endpoint: "NoRoute"}, defaultHandler, nil))
	engine.NoMethod(opencensus.HandlerFunc(&config.EndpointConfig{Endpoint: "NoMethod"}, defaultHandler, nil))

	logPrefix := "[SERVICE: Gin]"
	if err := httpsecure.Register(cfg.ExtraConfig, engine); err != nil && err != httpsecure.ErrNoConfig {
		opt.Logger.Warning(logPrefix+"[HTTPsecure]", err)
	} else if err == nil {
		opt.Logger.Debug(logPrefix + "[HTTPsecure] Successfuly loaded module")
	}

	lua.Register(opt.Logger, cfg.ExtraConfig, engine)

	botdetector.Register(cfg, opt.Logger, engine)

	return engine
}

func defaultHandler(c *gin.Context) {
	c.Header(core.KrakendHeaderName, core.KrakendHeaderValue)
	c.Header(server.CompleteResponseHeaderName, server.HeaderIncompleteResponseValue)
}

type engineFactory struct{}

func (engineFactory) NewEngine(cfg config.ServiceConfig, opt luragin.EngineOptions) *gin.Engine {
	return NewEngine(cfg, opt)
}
