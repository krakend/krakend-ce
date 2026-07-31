package krakend

import (
	"fmt"

	cel "github.com/krakend/krakend-cel/v3"
	jsonschema "github.com/krakend/krakend-jsonschema/v3"
	lua "github.com/krakend/krakend-lua/v3/proxy"
	metrics "github.com/krakend/krakend-metrics/v3/gin"
	"github.com/luraproject/lura/v3/config"
	"github.com/luraproject/lura/v3/logging"
	"github.com/luraproject/lura/v3/proxy"
)

func internalNewProxyFactory(logger logging.Logger, backendFactory proxy.BackendFactory,
	metricCollector *metrics.Metrics,
) proxy.Factory {
	proxyFactory := proxy.NewDefaultFactory(backendFactory, logger)
	proxyFactory = proxy.NewShadowFactory(proxyFactory)
	proxyFactory = jsonschema.ProxyFactory(logger, proxyFactory)
	proxyFactory = cel.ProxyFactory(logger, proxyFactory)
	proxyFactory = lua.ProxyFactory(logger, proxyFactory)
	proxyFactory = metricCollector.ProxyFactory("pipe", proxyFactory)
	return proxyFactory
}

// NewProxyFactory returns a new ProxyFactory wrapping the injected BackendFactory with the default proxy stack and a metrics collector
func NewProxyFactory(logger logging.Logger, backendFactory proxy.BackendFactory, metricCollector *metrics.Metrics) proxy.Factory {
	proxyFactory := internalNewProxyFactory(logger, backendFactory, metricCollector)

	return proxy.FactoryFunc(func(cfg *config.EndpointConfig) (proxy.Proxy, error) {
		logger.Debug(fmt.Sprintf("[ENDPOINT: %s] Building the proxy pipe", cfg.Endpoint))
		return proxyFactory.New(cfg)
	})
}

type proxyFactory struct{}

func (proxyFactory) NewProxyFactory(logger logging.Logger, backendFactory proxy.BackendFactory, metricCollector *metrics.Metrics) proxy.Factory {
	return NewProxyFactory(logger, backendFactory, metricCollector)
}
