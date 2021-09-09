package krakend

import (
	errorHandler "github.com/Unacademy/krakend-error-handler"
	cel "github.com/devopsfaith/krakend-cel"
	jsonschema "github.com/devopsfaith/krakend-jsonschema"
	lua "github.com/devopsfaith/krakend-lua/proxy"
	metrics "github.com/devopsfaith/krakend-metrics/gin"
	opencensus "github.com/devopsfaith/krakend-opencensus"
	"github.com/luraproject/lura/logging"
	"github.com/luraproject/lura/proxy"
)

// NewProxyFactory returns a new ProxyFactory wrapping the injected BackendFactory with the default proxy stack and a metrics collector
func NewProxyFactory(logger logging.Logger, backendFactory proxy.BackendFactory, metricCollector *metrics.Metrics) proxy.Factory {
	proxyFactory := proxy.NewDefaultFactory(backendFactory, logger)
	proxyFactory = proxy.NewShadowFactory(proxyFactory)
	proxyFactory = errorHandler.ProxyFactory(proxyFactory)
	proxyFactory = jsonschema.ProxyFactory(proxyFactory)
	proxyFactory = cel.ProxyFactory(logger, proxyFactory)
	proxyFactory = lua.ProxyFactory(logger, proxyFactory)
	proxyFactory = metricCollector.ProxyFactory("pipe", proxyFactory)
	proxyFactory = opencensus.ProxyFactory(proxyFactory)
	return proxyFactory
}

type proxyFactory struct{}

func (p proxyFactory) NewProxyFactory(logger logging.Logger, backendFactory proxy.BackendFactory, metricCollector *metrics.Metrics) proxy.Factory {
	return NewProxyFactory(logger, backendFactory, metricCollector)
}
