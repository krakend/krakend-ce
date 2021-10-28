package krakend

import (
	botdetector "github.com/devopsfaith/krakend-botdetector/v2/gin"
	jose "github.com/devopsfaith/krakend-jose/v2"
	ginjose "github.com/devopsfaith/krakend-jose/v2/gin"
	lua "github.com/devopsfaith/krakend-lua/v2/router/gin"
	metrics "github.com/devopsfaith/krakend-metrics/v2/gin"
	opencensus "github.com/devopsfaith/krakend-opencensus/v2/router/gin"
	juju "github.com/devopsfaith/krakend-ratelimit/v2/juju/router/gin"
	"github.com/luraproject/lura/v2/logging"
	router "github.com/luraproject/lura/v2/router/gin"
	"github.com/luraproject/lura/v2/transport/http/server"
)

// NewHandlerFactory returns a HandlerFactory with a rate-limit and a metrics collector middleware injected
func NewHandlerFactory(logger logging.Logger, metricCollector *metrics.Metrics, rejecter jose.RejecterFactory) router.HandlerFactory {
	handlerFactory := router.CustomErrorEndpointHandler(logger, server.DefaultToHTTPError)
	handlerFactory = juju.NewRateLimiterMw(handlerFactory)
	handlerFactory = lua.HandlerFactory(logger, handlerFactory)
	handlerFactory = ginjose.HandlerFactory(handlerFactory, logger, rejecter)
	handlerFactory = metricCollector.NewHTTPHandlerFactory(handlerFactory)
	handlerFactory = opencensus.New(handlerFactory)
	handlerFactory = botdetector.New(handlerFactory, logger)
	return handlerFactory
}

type handlerFactory struct{}

func (h handlerFactory) NewHandlerFactory(l logging.Logger, m *metrics.Metrics, r jose.RejecterFactory) router.HandlerFactory {
	return NewHandlerFactory(l, m, r)
}
