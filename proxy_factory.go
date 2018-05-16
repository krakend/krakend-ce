package main

import (
	metrics "github.com/devopsfaith/krakend-metrics/gin"
	"github.com/devopsfaith/krakend/config"
	"github.com/devopsfaith/krakend/logging"
	"github.com/devopsfaith/krakend/proxy"
)

// NewProxyFactory returns a new ProxyFactory wrapping the injected BackendFactory with the default proxy stack and a metrics collector
func NewProxyFactory(logger logging.Logger, backendFactory proxy.BackendFactory, metricCollector *metrics.Metrics) proxy.Factory {
	proxyFactory := proxy.NewDefaultFactory(backendFactory, logger)
	proxyFactory = NewLoggingProxyFactory(logger, proxyFactory)
	proxyFactory = metricCollector.ProxyFactory("pipe", proxyFactory)
	return proxyFactory
}

func NewLoggingProxyFactory(logger logging.Logger, proxyFactory proxy.Factory) proxy.FactoryFunc {
	return func(cfg *config.EndpointConfig) (proxy.Proxy, error) {
		next, err := proxyFactory.New(cfg)
		if err != nil {
			return next, err
		}
		return proxy.NewLoggingMiddleware(logger, "pipe")(next), nil
	}
}
