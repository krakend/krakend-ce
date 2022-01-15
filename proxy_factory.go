package krakend

import (
	cel "github.com/devopsfaith/krakend-cel/v2"
	jsonschema "github.com/devopsfaith/krakend-jsonschema/v2"
	lua "github.com/devopsfaith/krakend-lua/v2/proxy"
	metrics "github.com/devopsfaith/krakend-metrics/v2/gin"
	opencensus "github.com/devopsfaith/krakend-opencensus/v2"
	"github.com/luraproject/lura/v2/logging"
	"github.com/luraproject/lura/v2/proxy"
)

// NewProxyFactory returns a new ProxyFactory wrapping the injected BackendFactory with the default proxy stack and a metrics collector
func NewProxyFactory(logger logging.Logger, backendFactory proxy.BackendFactory, metricCollector *metrics.Metrics) proxy.Factory {
	proxyFactory := proxy.NewDefaultFactory(backendFactory, logger)
	proxyFactory = proxy.NewShadowFactory(proxyFactory)
	proxyFactory = jsonschema.ProxyFactory(logger, proxyFactory)
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
