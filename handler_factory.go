package krakend

import (
	"github.com/devopsfaith/krakend-jose"
	metrics "github.com/devopsfaith/krakend-metrics/gin"
	juju "github.com/devopsfaith/krakend-ratelimit/juju/router/gin"
	"github.com/devopsfaith/krakend/logging"
	router "github.com/devopsfaith/krakend/router/gin"
)

// NewHandlerFactory returns a HandlerFactory with a rate-limit and a metrics collector middleware injected
func NewHandlerFactory(logger logging.Logger, lcfg loggingConfig, metricCollector *metrics.Metrics, rejecter jose.RejecterFactory) router.HandlerFactory {
	router.RegisterRender("json_error", jsonErrorRender)
	handlerFactory := juju.HandlerFactory
	handlerFactory = NewJoseHandlerFactory(handlerFactory, logger, rejecter)
	handlerFactory = metricCollector.NewHTTPHandlerFactory(handlerFactory)
	return NewOpenCensusHandlerFactory(handlerFactory, lcfg)
}
