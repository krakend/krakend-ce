package krakend

import (
	"errors"
	"fmt"

	otelconfig "github.com/krakend/krakend-otel/config"
	otellura "github.com/krakend/krakend-otel/lura"
	otelstate "github.com/krakend/krakend-otel/state"
	cel "github.com/krakendio/krakend-cel/v2"
	jsonschema "github.com/krakendio/krakend-jsonschema/v2"
	lua "github.com/krakendio/krakend-lua/v2/proxy"
	metrics "github.com/krakendio/krakend-metrics/v2/gin"
	opencensus "github.com/krakendio/krakend-opencensus/v2"
	"github.com/luraproject/lura/v2/config"
	"github.com/luraproject/lura/v2/logging"
	"github.com/luraproject/lura/v2/proxy"
)

func newProxyFactoryWithConfig(logger logging.Logger, backendFactory proxy.BackendFactory,
	metricCollector *metrics.Metrics, serviceCfg *config.ServiceConfig) proxy.Factory {

	proxyFactory := proxy.NewDefaultFactory(backendFactory, logger)
	proxyFactory = proxy.NewShadowFactory(proxyFactory)
	proxyFactory = jsonschema.ProxyFactory(logger, proxyFactory)
	proxyFactory = cel.ProxyFactory(logger, proxyFactory)
	proxyFactory = lua.ProxyFactory(logger, proxyFactory)
	proxyFactory = metricCollector.ProxyFactory("pipe", proxyFactory)
	proxyFactory = opencensus.ProxyFactory(proxyFactory)

	if serviceCfg != nil {
		otelCfg, err := otelconfig.FromLura(*serviceCfg)
		if err != nil {
			if !errors.Is(err, otelconfig.ErrNoConfig) {
				logger.Error(fmt.Sprintf("cannot load OpenTelemetry config: %s", err.Error()))
			}
		} else {
			proxyFactory = otellura.ProxyFactory(proxyFactory, otelstate.GlobalState, otelCfg.Layers.Pipe,
				otelCfg.SkipPaths)
		}
	}

	return proxyFactory
}

// NewProxyFactory returns a new ProxyFactory wrapping the injected BackendFactory with the default proxy stack and a metrics collector
func NewProxyFactory(logger logging.Logger, backendFactory proxy.BackendFactory, metricCollector *metrics.Metrics) proxy.Factory {
	proxyFactory := newProxyFactoryWithConfig(logger, backendFactory, metricCollector, nil)

	return proxy.FactoryFunc(func(cfg *config.EndpointConfig) (proxy.Proxy, error) {
		logger.Debug(fmt.Sprintf("[ENDPOINT: %s] Building the proxy pipe", cfg.Endpoint))
		return proxyFactory.New(cfg)
	})
}

func newProxyFactoryWithConfigAndDbg(logger logging.Logger, backendFactory proxy.BackendFactory, metricCollector *metrics.Metrics,
	serviceCfg *config.ServiceConfig) proxy.Factory {

	proxyFactory := newProxyFactoryWithConfig(logger, backendFactory, metricCollector, serviceCfg)

	return proxy.FactoryFunc(func(cfg *config.EndpointConfig) (proxy.Proxy, error) {
		logger.Debug(fmt.Sprintf("[ENDPOINT: %s] Building the proxy pipe", cfg.Endpoint))
		return proxyFactory.New(cfg)
	})
}

type proxyFactory struct{}

func (proxyFactory) NewProxyFactory(logger logging.Logger, backendFactory proxy.BackendFactory, metricCollector *metrics.Metrics) proxy.Factory {
	return NewProxyFactory(logger, backendFactory, metricCollector)
}

func (proxyFactory) NewProxyFactoryWithConfig(logger logging.Logger, backendFactory proxy.BackendFactory, metricCollector *metrics.Metrics,
	serviceCfg *config.ServiceConfig) proxy.Factory {

	return newProxyFactoryWithConfigAndDbg(logger, backendFactory, metricCollector, serviceCfg)
}
