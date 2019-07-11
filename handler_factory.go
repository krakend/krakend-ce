package krakend

import (
	botdetector "github.com/devopsfaith/krakend-botdetector/gin"
	"github.com/devopsfaith/krakend-jose"
	ginjose "github.com/devopsfaith/krakend-jose/gin"
	lua "github.com/devopsfaith/krakend-lua/router/gin"
	metrics "github.com/devopsfaith/krakend-metrics/gin"
	opencensus "github.com/devopsfaith/krakend-opencensus/router/gin"
	juju "github.com/devopsfaith/krakend-ratelimit/juju/router/gin"
	"github.com/devopsfaith/krakend/logging"
	router "github.com/devopsfaith/krakend/router/gin"
)

// NewHandlerFactory returns a HandlerFactory with a rate-limit and a metrics collector middleware injected
func NewHandlerFactory(logger logging.Logger, metricCollector *metrics.Metrics, rejecter jose.RejecterFactory) router.HandlerFactory {
	handlerFactory := juju.HandlerFactory
	handlerFactory = lua.HandlerFactory(logger, handlerFactory)
	handlerFactory = ginjose.HandlerFactory(handlerFactory, logger, rejecter)
	handlerFactory = metricCollector.NewHTTPHandlerFactory(handlerFactory)
	handlerFactory = opencensus.New(handlerFactory)
	handlerFactory = botdetector.New(handlerFactory, logger)
	return handlerFactory
}
