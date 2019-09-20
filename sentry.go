package krakend

import (
	"strconv"
    "net/http"
	"github.com/devopsfaith/krakend/config"
	"go.opencensus.io/trace"
	"github.com/getsentry/sentry-go"
	sentrygin "github.com/getsentry/sentry-go/gin"
	"github.com/gin-gonic/gin"
)

const sentryNamespace = "github_com/openrm/sentry"

var nopHandler = func(c *gin.Context) { c.Next() }

func NewSentryMiddleware(cfg config.ServiceConfig) gin.HandlerFunc {
	data, ok := cfg.ExtraConfig[sentryNamespace]

	if !ok {
		return nopHandler
	}

	var dsn, environment string

	if cfg, ok := data.(map[string]interface{}); ok {
		if v, ok := cfg["dsn"]; ok {
			if s, ok := v.(string); ok {
				dsn = s
			}
		}
		if v, ok := cfg["environment"]; ok {
			if s, ok := v.(string); ok {
				environment = s
			}
		}
	}

	if err := sentry.Init(sentry.ClientOptions{
		Dsn: dsn,
		Debug: true,
		Release: strconv.Itoa(cfg.Version),
		Environment: environment,
	}); err != nil {
		return nopHandler
	}

	handler := sentrygin.New(sentrygin.Options{
		Repanic: true,
	})

	return func(c *gin.Context) {
		handler(c)
        if len(c.Errors) > 0 {
            if status := c.Writer.Status(); status >= http.StatusInternalServerError {
                if hub := sentrygin.GetHubFromContext(c); hub != nil {
                    if span := trace.FromContext(c.Request.Context()); span != nil {
                        sc := span.SpanContext()
                        hub.Scope().SetContext("trace", map[string]interface{}{
                            "trace_id": sc.TraceID.String(),
                            "span_id": sc.SpanID.String(),
                        })
                    }
                    for _, err := range c.Errors {
                        hub.CaptureException(err)
                    }
                }
            }
		}
	}
}
