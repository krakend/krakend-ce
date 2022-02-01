package krakend

import (
	"io"
	"time"

	botdetector "github.com/devopsfaith/krakend-botdetector/gin"
	httpsecure "github.com/devopsfaith/krakend-httpsecure/gin"
	lua "github.com/devopsfaith/krakend-lua/router/gin"
	"github.com/gin-gonic/gin"
	"github.com/luraproject/lura/config"
	"github.com/luraproject/lura/logging"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// NewEngine creates a new gin engine with some default values and a secure middleware
func NewEngine(cfg config.ServiceConfig, logger logging.Logger, w io.Writer) *gin.Engine {
	if !cfg.Debug {
		gin.SetMode(gin.ReleaseMode)
	}

	engine := gin.New()
	engine.Use(zapLogger(), gin.Recovery())

	engine.RedirectTrailingSlash = true
	engine.RedirectFixedPath = true
	engine.HandleMethodNotAllowed = true

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

// zapLogger instance a Logger middleware with config.
func zapLogger() gin.HandlerFunc {
	logger, _ := zap.NewProduction()

	// TODO: populate with skipped paths via config
	var skip map[string]struct{}

	return func(c *gin.Context) {
		// Start timer
		// TODO: add utc support
		start := time.Now()
		path := c.Request.URL.Path

		// Process request
		c.Next()

		// Log only when path is not being skipped
		if _, ok := skip[path]; !ok {
			// TODO: add utc support
			end := time.Now()
			status := c.Writer.Status()
			latency := end.Sub(start)

			fields := []zapcore.Field{
				zap.Duration("duration", latency),
				zap.String("host", c.Request.URL.Host),
				zap.String("ip", c.ClientIP()),
				zap.String("method", c.Request.Method),
				zap.String("path", path),
				zap.Int("status", status),
				zap.String("user_agent", c.Request.UserAgent()),
			}
			if len(c.Errors) > 0 {
				fields = append(fields, zap.Strings("errors", c.Errors.Errors()))
			}

			if status >= 200 && status < 400 {
				logger.Info("request", fields...)
			} else if status >= 400 && status < 500 {
				logger.Warn("bad request", fields...)
			} else {
				logger.Error("request error", fields...)
			}
		}
	}
}
