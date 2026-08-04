package krakend

import (
	"context"
	"fmt"

	amqp "github.com/krakend/krakend-amqp/v3"
	cel "github.com/krakend/krakend-cel/v3"
	cb "github.com/krakend/krakend-circuitbreaker/v4/gobreaker/proxy"
	httpcache "github.com/krakend/krakend-httpcache/v3"
	lambda "github.com/krakend/krakend-lambda/v3"
	lua "github.com/krakend/krakend-lua/v3/proxy"
	martian "github.com/krakend/krakend-martian/v3"
	metrics "github.com/krakend/krakend-metrics/v3/gin"
	oauth2client "github.com/krakend/krakend-oauth2-clientcredentials/v3"
	otellura "github.com/krakend/krakend-otel/v2/lura"
	pubsub "github.com/krakend/krakend-pubsub/v3"
	ratelimit "github.com/krakend/krakend-ratelimit/v4/proxy"
	"github.com/luraproject/lura/v3/config"
	"github.com/luraproject/lura/v3/logging"
	"github.com/luraproject/lura/v3/proxy"
	"github.com/luraproject/lura/v3/transport/http/client"
	httprequestexecutor "github.com/luraproject/lura/v3/transport/http/client/plugin"
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
func NewBackendFactory(logger logging.Logger, metricCollector *metrics.Metrics) proxy.BackendFactory {
	return NewBackendFactoryWithContext(context.Background(), logger, metricCollector)
}

func newRequestExecutorFactory(ctx context.Context, logger logging.Logger) func(*config.Backend) client.HTTPRequestExecutor {
	requestExecutorFactory := func(cfg *config.Backend) client.HTTPRequestExecutor {
		clientFactory := client.NewHTTPClient
		if _, ok := cfg.ExtraConfig[oauth2client.Namespace]; ok {
			clientFactory = oauth2client.NewHTTPClient(cfg)
		}

		clientFactory = httpcache.NewHTTPClient(cfg, clientFactory)
		clientFactory = otellura.InstrumentedHTTPClientFactory(clientFactory, cfg)
		return client.DefaultHTTPRequestExecutor(clientFactory)
	}
	return httprequestexecutor.HTTPRequestExecutorWithContext(ctx, logger, requestExecutorFactory)
}

func internalNewBackendFactory(
	ctx context.Context,
	requestExecutorFactory func(*config.Backend) client.HTTPRequestExecutor,
	logger logging.Logger,
	metricCollector *metrics.Metrics,
) proxy.BackendFactory {
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
	backendFactory = otellura.BackendFactory(backendFactory)
	return func(remote *config.Backend) proxy.Proxy {
		logger.Debug(fmt.Sprintf("[BACKEND: %s] Building the backend pipe", remote.URLPattern))
		return backendFactory(remote)
	}
}

// NewBackendFactoryWithContext creates a BackendFactory by stacking all the available middlewares and injecting the received context
func NewBackendFactoryWithContext(ctx context.Context, logger logging.Logger, metricCollector *metrics.Metrics) proxy.BackendFactory {
	requestExecutorFactory := newRequestExecutorFactory(ctx, logger)
	return internalNewBackendFactory(ctx, requestExecutorFactory, logger, metricCollector)
}

type backendFactory struct{}

func (backendFactory) NewBackendFactory(ctx context.Context, l logging.Logger, m *metrics.Metrics) proxy.BackendFactory {
	return NewBackendFactoryWithContext(ctx, l, m)
}
