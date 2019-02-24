package krakend

import (
	cel "github.com/devopsfaith/krakend-cel"
	cb "github.com/devopsfaith/krakend-circuitbreaker/gobreaker/proxy"
	httpcache "github.com/devopsfaith/krakend-httpcache"
	"github.com/devopsfaith/krakend-martian"
	metrics "github.com/devopsfaith/krakend-metrics/gin"
	"github.com/devopsfaith/krakend-oauth2-clientcredentials"
	opencensus "github.com/devopsfaith/krakend-opencensus"
	juju "github.com/devopsfaith/krakend-ratelimit/juju/proxy"
	"github.com/devopsfaith/krakend/config"
	"github.com/devopsfaith/krakend/logging"
	"github.com/devopsfaith/krakend/proxy"
	"github.com/devopsfaith/krakend/transport/http/client"
)

// NewBackendFactory creates a BackendFactory by stacking all the available middlewares:
// - oauth2 client credentials
// - martian
// - rate-limit
// - circuit breaker
// - metrics collector
// - opencensus collector
func NewBackendFactory(logger logging.Logger, metricCollector *metrics.Metrics) proxy.BackendFactory {
	requestExecutorFactory := func(cfg *config.Backend) client.HTTPRequestExecutor {
		var clientFactory client.HTTPClientFactory
		if _, ok := cfg.ExtraConfig[oauth2client.Namespace]; ok {
			clientFactory = oauth2client.NewHTTPClient(cfg)
		} else {
			clientFactory = httpcache.NewHTTPClient(cfg)
		}
		return opencensus.HTTPRequestExecutor(clientFactory)
	}
	backendFactory := martian.NewConfiguredBackendFactory(logger, requestExecutorFactory)
	backendFactory = cel.BackendFactory(logger, backendFactory)
	backendFactory = juju.BackendFactory(backendFactory)
	backendFactory = cb.BackendFactory(backendFactory, logger)
	backendFactory = metricCollector.BackendFactory("backend", backendFactory)
	backendFactory = opencensus.BackendFactory(backendFactory)
	return backendFactory
}
