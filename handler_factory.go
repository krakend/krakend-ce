package main

import (
	metrics "github.com/devopsfaith/krakend-metrics/gin"
	juju "github.com/devopsfaith/krakend-ratelimit/juju/router/gin"
	"github.com/devopsfaith/krakend/logging"
	router "github.com/devopsfaith/krakend/router/gin"
)

// NewHandlerFactory returns a HandlerFactory with a rate-limit and a metrics collector middleware injected
func NewHandlerFactory(logger logging.Logger, metricCollector *metrics.Metrics) router.HandlerFactory {
	handlerFactory := juju.HandlerFactory
	handlerFactory = metricCollector.NewHTTPHandlerFactory(handlerFactory)
	return handlerFactory
}
