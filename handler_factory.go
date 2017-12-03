package main

import (
	metrics "github.com/devopsfaith/krakend-metrics/gin"
	juju "github.com/devopsfaith/krakend-ratelimit/juju/router/gin"
	router "github.com/devopsfaith/krakend/router/gin"
)

// NewHandlerFactory returns a HandlerFactory with a rate-limit and a metrics collector middleware injected
func NewHandlerFactory(metricCollector *metrics.Metrics) router.HandlerFactory {
	handlerFactory := juju.HandlerFactory
	handlerFactory = metricCollector.NewHTTPHandlerFactory(handlerFactory)
	return handlerFactory
}
