package krakend

import (
	"context"
	"fmt"

	otellura "github.com/krakend/krakend-otel/lura"
	amqp "github.com/krakendio/krakend-amqp/v2"
	cel "github.com/krakendio/krakend-cel/v2"
	cb "github.com/krakendio/krakend-circuitbreaker/v2/gobreaker/proxy"
	httpcache "github.com/krakendio/krakend-httpcache/v2"
	lambda "github.com/krakendio/krakend-lambda/v2"
	lua "github.com/krakendio/krakend-lua/v2/proxy"
	martian "github.com/krakendio/krakend-martian/v2"
	metrics "github.com/krakendio/krakend-metrics/v2/gin"
	oauth2client "github.com/krakendio/krakend-oauth2-clientcredentials/v2"
	opencensus "github.com/krakendio/krakend-opencensus/v2"
	pubsub "github.com/krakendio/krakend-pubsub/v2"
	ratelimit "github.com/krakendio/krakend-ratelimit/v3/proxy"
	"github.com/luraproject/lura/v2/config"
	"github.com/luraproject/lura/v2/logging"
	"github.com/luraproject/lura/v2/proxy"
	"github.com/luraproject/lura/v2/transport/http/client"
	httprequestexecutor "github.com/luraproject/lura/v2/transport/http/client/plugin"
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

func newRequestExecutorFactory(logger logging.Logger) func(*config.Backend) client.HTTPRequestExecutor {
	requestExecutorFactory := func(cfg *config.Backend) client.HTTPRequestExecutor {
		clientFactory := client.NewHTTPClient
		if _, ok := cfg.ExtraConfig[oauth2client.Namespace]; ok {
			clientFactory = oauth2client.NewHTTPClient(cfg)
		}

		clientFactory = httpcache.NewHTTPClient(cfg, clientFactory)
		clientFactory = otellura.InstrumentedHTTPClientFactory(clientFactory, cfg)
		// TODO: check what happens if we have both, opencensus and otel enabled ?
		return opencensus.HTTPRequestExecutorFromConfig(clientFactory, cfg)
	}
	return httprequestexecutor.HTTPRequestExecutor(logger, requestExecutorFactory)
}

func internalNewBackendFactory(ctx context.Context, requestExecutorFactory func(*config.Backend) client.HTTPRequestExecutor,
	logger logging.Logger, metricCollector *metrics.Metrics) proxy.BackendFactory {

	backendFactory := martian.NewConfiguredBackendFactory(logger, requestExecutorFactory)
	bf := pubsub.NewBackendFactory(ctx, logger, backendFactory)
	backendFactory = bf.New
	backendFactory = amqp.NewBackendFactory(ctx, logger, backendFactory)
	backendFactory = lambda.BackendFactory(logger, backendFactory)
	backendFactory = cel.BackendFactory(logger, backendFactory)
	backendFactory = lua.BackendFactory(logger, backendFactory)
	backendFactory = ratelimit.BackendFactory(logger, backendFactory)
	backendFactory = cb.BackendFactory(backendFactory, logger)
	backendFactory = metricCollector.BackendFactory("backend", backendFactory)
	backendFactory = opencensus.BackendFactory(backendFactory)
	backendFactory = otellura.BackendFactory(backendFactory)
	return func(remote *config.Backend) proxy.Proxy {
		logger.Debug(fmt.Sprintf("[BACKEND: %s] Building the backend pipe", remote.URLPattern))
		return backendFactory(remote)
	}
}

// NewBackendFactoryWithContext creates a BackendFactory by stacking all the available middlewares and injecting the received context
func NewBackendFactoryWithContext(ctx context.Context, logger logging.Logger, metricCollector *metrics.Metrics) proxy.BackendFactory {
	requestExecutorFactory := newRequestExecutorFactory(logger)
	return internalNewBackendFactory(ctx, requestExecutorFactory, logger, metricCollector)
}

type backendFactory struct{}

func (backendFactory) NewBackendFactory(ctx context.Context, l logging.Logger, m *metrics.Metrics) proxy.BackendFactory {
	return NewBackendFactoryWithContext(ctx, l, m)
}
