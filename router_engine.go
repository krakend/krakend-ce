package krakend

import (
	"io"

	botdetector "github.com/devopsfaith/krakend-botdetector/v2/gin"
	httpsecure "github.com/devopsfaith/krakend-httpsecure/v2/gin"
	lua "github.com/devopsfaith/krakend-lua/v2/router/gin"
	"github.com/gin-gonic/gin"
	"github.com/luraproject/lura/v2/config"
	"github.com/luraproject/lura/v2/logging"
	luragin "github.com/luraproject/lura/v2/router/gin"
)

// NewEngine creates a new gin engine with some default values and a secure middleware
func NewEngine(cfg config.ServiceConfig, logger logging.Logger, w io.Writer) *gin.Engine {
	engine := luragin.NewEngine(cfg, logger, w)

	if err := httpsecure.Register(cfg.ExtraConfig, engine); err != nil {
		logger.Warning(err)
	}

	lua.Register(logger, cfg.ExtraConfig, engine)

	botdetector.Register(cfg, logger, engine)

	return engine
}

type engineFactory struct{}

func (e engineFactory) NewEngine(cfg config.ServiceConfig, l logging.Logger, w io.Writer) *gin.Engine {
	return NewEngine(cfg, l, w)
}
