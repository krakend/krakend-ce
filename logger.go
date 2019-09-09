package krakend

import (
	"os"
	"time"
	"github.com/sirupsen/logrus"
	"github.com/gin-gonic/gin"
	"github.com/devopsfaith/krakend/config"
	tracing "github.com/openrm/module-tracing-golang/opentracing"
)

const LoggerNamespace = "github_com/openrm/logging"

func newLoggingHandler(logger logrus.FieldLogger, skip map[string]struct{}, traceHeader string) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		r := c.Request
		headers := r.Header

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

		if h := headers.Get(traceHeader); h != "" {
			if span := tracing.NewFromTraceParent(h); span != nil {
				entry = entry.WithFields(logrus.Fields{
					"span": span.JSON(true),
					"traceId": span.TraceId,
					"spanId": span.SpanId,
				})
			}
		}

		c.Next()

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
				if ps, ok := v.([]string); ok && len(ps) > 0 {
					skip = make(map[string]struct{}, len(ps))
					for _, path := range ps {
						skip[path] = struct{}{}
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

	return newLoggingHandler(logger, skip, traceHeader)
}
