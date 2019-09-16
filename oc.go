package krakend

import (
	"context"
	"net/http"
	"github.com/gin-gonic/gin"
	"github.com/devopsfaith/krakend/transport/http/client"
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
		client.Transport = &ochttp.Transport{
			Base: transport,
			Propagation: lcfg.httpFormat(),
		}
		return client
	}
}

func NewOpenCensusMiddleware(lcfg loggingConfig) gin.HandlerFunc {
	skip := lcfg.skipPaths
	handler := ochttp.Handler{
		Propagation: lcfg.httpFormat(),
		GetStartOptions: func(r *http.Request) trace.StartOptions{
			if u := r.URL; u != nil {
				if _, ok := skip[u.Path]; ok {
					return trace.StartOptions{
						Sampler: trace.NeverSample(),
					}
				}
			}
			return trace.StartOptions{}
		},
	}
	return func(c *gin.Context) {
		handler.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			c.Request = r
			c.Next()
		})
		handler.ServeHTTP(c.Writer, c.Request)
	}
}
