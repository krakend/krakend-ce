package krakend

import (
	botdetector "github.com/devopsfaith/krakend-botdetector/v2/gin"
	httpsecure "github.com/devopsfaith/krakend-httpsecure/v2/gin"
	lua "github.com/devopsfaith/krakend-lua/v2/router/gin"
	"github.com/gin-gonic/gin"
	"github.com/luraproject/lura/v2/config"
	luragin "github.com/luraproject/lura/v2/router/gin"
)

// NewEngine creates a new gin engine with some default values and a secure middleware
func NewEngine(cfg config.ServiceConfig, opt luragin.EngineOptions) *gin.Engine {
	engine := luragin.NewEngine(cfg, opt)
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

type engineFactory struct{}

func (e engineFactory) NewEngine(cfg config.ServiceConfig, opt luragin.EngineOptions) *gin.Engine {
	return NewEngine(cfg, opt)
}
