package krakend

import (
	"context"

	amqp "github.com/devopsfaith/krakend-amqp"
	cel "github.com/devopsfaith/krakend-cel"
	cb "github.com/devopsfaith/krakend-circuitbreaker/gobreaker/proxy"
	httpcache "github.com/devopsfaith/krakend-httpcache"
	lambda "github.com/devopsfaith/krakend-lambda"
	lua "github.com/devopsfaith/krakend-lua/proxy"
	"github.com/devopsfaith/krakend-martian"
	metrics "github.com/devopsfaith/krakend-metrics/gin"
	"github.com/devopsfaith/krakend-oauth2-clientcredentials"
	opencensus "github.com/devopsfaith/krakend-opencensus"
	pubsub "github.com/devopsfaith/krakend-pubsub"
	juju "github.com/devopsfaith/krakend-ratelimit/juju/proxy"
	"github.com/luraproject/lura/config"
	"github.com/luraproject/lura/logging"
	"github.com/luraproject/lura/proxy"
	"github.com/luraproject/lura/transport/http/client"
	httprequestexecutor "github.com/luraproject/lura/transport/http/client/plugin"
)

// NewBackendFactory creates a BackendFactory by stacking all the available middlewares:
// - oauth2 client credentials
// - http cache
// - martian
// - pubsub
// - amqp
// - cel
// - lua
// - rate-limit
// - circuit breaker
// - metrics collector
// - opencensus collector
func NewBackendFactory(logger logging.Logger, metricCollector *metrics.Metrics) proxy.BackendFactory {
	return NewBackendFactoryWithContext(context.Background(), logger, metricCollector)
}

// NewBackendFactory creates a BackendFactory by stacking all the available middlewares and injecting the received context
func NewBackendFactoryWithContext(ctx context.Context, logger logging.Logger, metricCollector *metrics.Metrics) proxy.BackendFactory {
	requestExecutorFactory := func(cfg *config.Backend) client.HTTPRequestExecutor {
		var clientFactory client.HTTPClientFactory
		if _, ok := cfg.ExtraConfig[oauth2client.Namespace]; ok {
			clientFactory = oauth2client.NewHTTPClient(cfg)
		} else {
			clientFactory = httpcache.NewHTTPClient(cfg)
		}
		return opencensus.HTTPRequestExecutor(clientFactory)
	}
	requestExecutorFactory = httprequestexecutor.HTTPRequestExecutor(logger, requestExecutorFactory)
	backendFactory := martian.NewConfiguredBackendFactory(logger, requestExecutorFactory)
	bf := pubsub.NewBackendFactory(ctx, logger, backendFactory)
	backendFactory = bf.New
	backendFactory = amqp.NewBackendFactory(ctx, logger, backendFactory)
	backendFactory = lambda.BackendFactory(backendFactory)
	backendFactory = cel.BackendFactory(logger, backendFactory)
	backendFactory = lua.BackendFactory(logger, backendFactory)
	backendFactory = juju.BackendFactory(backendFactory)
	backendFactory = cb.BackendFactory(backendFactory, logger)
	backendFactory = metricCollector.BackendFactory("backend", backendFactory)
	backendFactory = opencensus.BackendFactory(backendFactory)
	return backendFactory
}

type backendFactory struct{}

func (b backendFactory) NewBackendFactory(ctx context.Context, l logging.Logger, m *metrics.Metrics) proxy.BackendFactory {
	return NewBackendFactoryWithContext(ctx, l, m)
}
