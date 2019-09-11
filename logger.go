package krakend

import (
	"os"
	"time"
	"github.com/sirupsen/logrus"
	"github.com/gin-gonic/gin"
	"github.com/devopsfaith/krakend/config"
	"go.opencensus.io/trace"
	"github.com/openrm/module-tracing-golang/propagation"
)

const LoggerNamespace = "github_com/openrm/logging"

func spanContextMap(sc trace.SpanContext) map[string]interface{} {
	return map[string]interface{}{
		"traceId": sc.TraceID.String(),
		"spanId": sc.SpanID.String(),
		"sampled": sc.IsSampled(),
	}
}

func loggingHandler(logger logrus.FieldLogger, skip map[string]struct{}, traceHeader string) gin.HandlerFunc {
	format := propagation.HTTPFormat{Header: traceHeader}

	return func(c *gin.Context) {
		start := time.Now()

		r := c.Request

		if u := r.URL; u != nil {
			if _, ok := skip[u.Path]; ok {
				c.Next()
				return
			}
		}

		entry := logger.WithField("module", "GIN")

		entry = entry.WithFields(logrus.Fields{
			"ip": c.ClientIP(),
			"method": r.Method,
			"protocol": r.Proto,
			"url": r.RequestURI,
			"remoteAddress": r.RemoteAddr,
			"hostname": r.Host,
			"referer": r.Referer(),
			"userAgent": r.UserAgent(),
			"contentLength": r.ContentLength,
		})

		traceData := make(map[string]interface{})
		if sc, ok := format.SpanContextFromRequest(r); ok{
			traceData["parent"] = spanContextMap(sc)
		}

		c.Next()

		if span := trace.FromContext(c.Request.Context()); span != nil {
			for k, v := range spanContextMap(span.SpanContext()) {
				traceData[k] = v
			}
			entry = entry.WithField("span", traceData)
		}

		w := c.Writer

		if len(c.Errors) > 0 {
			entry = entry.WithField("err", c.Errors.JSON())
		}

		entry.WithFields(logrus.Fields{
			"responseTime": float64(time.Now().Sub(start).Nanoseconds()) / 1e6,
			"status": w.Status(),
			"responseContentLength": w.Size(),
		}).Info()
	}
}

func NewRouterLogger(cfg config.ExtraConfig) gin.HandlerFunc {
	var skip map[string]struct{}
	var traceHeader string

	if e, ok := cfg[LoggerNamespace]; ok {
		if m, ok := e.(map[string]interface{}); ok {
			if v, ok := m["skip_paths"]; ok {
				if ps, ok := v.([]interface{}); ok && len(ps) > 0 {
					skip = make(map[string]struct{}, len(ps))
					for _, v := range ps {
						if path, ok := v.(string); ok {
							skip[path] = struct{}{}
						}
					}
				}
			}
			if v, ok := m["trace_header"]; ok {
				if h, ok := v.(string); ok {
					traceHeader = h
				}
			}
		}
	}

	logger := logrus.New()

	logger.SetOutput(os.Stdout)
	logger.SetFormatter(&logrus.JSONFormatter{
		FieldMap: logrus.FieldMap{
			logrus.FieldKeyTime: "@timestamp",
			logrus.FieldKeyMsg: "message",
		},
	})

	return loggingHandler(logger, skip, traceHeader)
}
