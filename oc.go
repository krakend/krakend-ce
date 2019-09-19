package krakend

import (
	"context"
	"net/http"
	"github.com/devopsfaith/krakend/config"
	"github.com/devopsfaith/krakend/proxy"
	"github.com/gin-gonic/gin"
	router "github.com/devopsfaith/krakend/router/gin"
	"github.com/devopsfaith/krakend/transport/http/client"
	"github.com/devopsfaith/krakend-opencensus"
	"go.opencensus.io/plugin/ochttp"
	"go.opencensus.io/trace"
)

func NewOpenCensusClient(lcfg loggingConfig, clientFactory client.HTTPClientFactory) client.HTTPClientFactory {
	return func(ctx context.Context) *http.Client {
		client := clientFactory(ctx)
		transport := client.Transport
		if transport == nil {
			transport = http.DefaultTransport
		}
		if _, ok := transport.(*ochttp.Transport); !ok {
			client.Transport = &ochttp.Transport{
				Base: transport,
				Propagation: lcfg.httpFormat(),
			}
		}
		return client
	}
}

func NewOpenCensusHandlerFactory(hf router.HandlerFactory, lcfg loggingConfig) router.HandlerFactory {
	skip, prop := lcfg.skipPaths, lcfg.httpFormat()
	filterPath := func(r *http.Request) trace.StartOptions{
		if u := r.URL; u != nil {
			if _, ok := skip[u.Path]; ok {
				return trace.StartOptions{
					Sampler: trace.NeverSample(),
				}
			}
		}
		return trace.StartOptions{}
	}
	return func(cfg *config.EndpointConfig, p proxy.Proxy) gin.HandlerFunc {
		handler := hf(cfg, p)
		traceHandler := ochttp.Handler{
			Propagation: prop,
			GetStartOptions: filterPath,
			FormatSpanName: func(*http.Request) string {
				return cfg.Endpoint
			},
		}
		return func(c *gin.Context) {
			traceHandler.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				c.Request = r
				c.Set(opencensus.ContextKey, trace.FromContext(c.Request.Context()))
				handler(c)
			})
			traceHandler.ServeHTTP(c.Writer, c.Request)
		}
	}
}
