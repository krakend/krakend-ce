package krakend

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	opencensus "github.com/krakendio/krakend-opencensus/v2"
	"github.com/luraproject/lura/v2/config"
	"github.com/luraproject/lura/v2/proxy"
	krakendgin "github.com/luraproject/lura/v2/router/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.opentelemetry.io/otel/attribute"
	prometheus "go.opentelemetry.io/otel/exporters/prometheus"
	"go.opentelemetry.io/otel/metric/instrument"
	metric "go.opentelemetry.io/otel/sdk/metric"
)

const (
	Namespace = "telemetry/otel"
)

var (
	ErrNoConfig     = errors.New("No open telemetry config found")
	InvalidExporter = errors.New("Invalid exporter found")
)

type OpenTelemetryMetrics struct {
	Provider *metric.MeterProvider
}

func InitializePrometheus(ctx context.Context) (*prometheus.Exporter, starter) {
	exporter, err := prometheus.New()
	if err != nil {
		log.Fatal("Error while initializing prometheus exporter", err)
	}
	return exporter, func(cfg *Config) {
		go serveMetrics(cfg)
	}
}

type starter func(*Config)

func getProvider(ctx context.Context, cfg *Config) *metric.MeterProvider {
	var exporter *prometheus.Exporter
	var start starter
	if cfg.Exporters.Prometheus != nil {
		exporter, start = InitializePrometheus(ctx)
		start(cfg)
	}
	if exporter == nil {
		log.Fatal("No exporter found")
	}
	return metric.NewMeterProvider(metric.WithReader(exporter))
}

func InitializeOpenTelemetryMetricsCollector(ctx context.Context, cfg config.ServiceConfig) *OpenTelemetryMetrics {
	extraConfig, ok := cfg.ExtraConfig[Namespace]
	if !ok {
		return nil
	}
	config, err := parseCfg(new(Config), extraConfig)
	if err != nil {
		log.Fatal("Unable to load config for open telemetry")
	}
	provider := getProvider(ctx, config)
	return &OpenTelemetryMetrics{Provider: provider}

	// gauge, err := meter.SyncFloat64().UpDownCounter("bar", instrument.WithDescription("a fun little gauge"))
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// gauge.Add(ctx, 100, attrs...)
	// gauge.Add(ctx, -25, attrs...)

	// // This is the equivalent of prometheus.NewHistogramVec
	// histogram, err := meter.SyncFloat64().Histogram("baz", instrument.WithDescription("a very nice histogram"))
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// histogram.Record(ctx, 23, attrs...)
	// histogram.Record(ctx, 7, attrs...)
	// histogram.Record(ctx, 101, attrs...)
	// histogram.Record(ctx, 105, attrs...)

	// ctx, _ = signal.NotifyContext(ctx, os.Interrupt)
	// <-ctx.Done()
}

func exitIfError(e error, msg string) {
	if e != nil {
		log.Fatalf(msg, e)
	}
}

func (otelMetrics *OpenTelemetryMetrics) HandlerFactory(hf krakendgin.HandlerFactory) krakendgin.HandlerFactory {
	if otelMetrics == nil {
		return hf
	}

	meter := otelMetrics.Provider.Meter("krakend.otel.http.server")

	// Start the prometheus HTTP server and pass the exporter Collector to it

	// This is the equivalent of prometheus.NewCounterVec
	counter, err := meter.SyncInt64().Counter("http.server.request.count", instrument.WithDescription("a request counter"))
	if err != nil {
		log.Fatal(err)
	}

	return func(cfg *config.EndpointConfig, p proxy.Proxy) gin.HandlerFunc {
		extraConfig, ok := cfg.ExtraConfig[Namespace]
		var tags Tags
		if ok {
			otelEndpointConfig, err := parseCfg(new(EndpointConfig), extraConfig)
			exitIfError(err, "Error occured while loading endpoing telemetry config")
			tags = otelEndpointConfig.Tags
		}
		fixedAttrs := []attribute.KeyValue{}
		for k, v := range tags {
			fixedAttrs = append(fixedAttrs, attribute.Key(k).String(v))
		}
		extractPath := opencensus.GetAggregatedPathForMetrics(cfg)
		next := hf(cfg, p)
		return func(c *gin.Context) {
			requestPath := extractPath(c.Request)
			next(c)
			responseStatus := c.Writer.Status()

			attrs := []attribute.KeyValue{
				attribute.Key("http_method").String(c.Request.Method),
				attribute.Key("http_status").Int(responseStatus),
				attribute.Key("http_path").String(requestPath),
			}
			finalAttrs := append(fixedAttrs, attrs...)
			counter.Add(c, 1, finalAttrs...)
		}
	}
}

type Tags map[string]string

type Config struct {
	Exporters Exporters `json:"exporters"`
	Tags      Tags      `json:"tags"`
}

type EndpointConfig struct {
	Tags Tags `json:"tags"`
}

type Exporters struct {
	Prometheus *PrometheusConfig `json:"prometheus"`
}

type PrometheusConfig struct {
	Port     int    `json:"port`
	Endpoint string `json:"endpoint"`
}

func parseCfg[T any](cfg *T, extraConfig interface{}) (*T, error) {
	buf := new(bytes.Buffer)
	json.NewEncoder(buf).Encode(extraConfig)
	if err := json.NewDecoder(buf).Decode(cfg); err != nil {
		return cfg, err
	}
	return cfg, nil
}

func serveMetrics(cfg *Config) {
	url := fmt.Sprintf("http://localhost:%d%s", cfg.Exporters.Prometheus.Port, cfg.Exporters.Prometheus.Endpoint)
	log.Printf("Serving metrics at %s", url)
	http.Handle(cfg.Exporters.Prometheus.Endpoint, promhttp.Handler())
	err := http.ListenAndServe(fmt.Sprintf(":%d", cfg.Exporters.Prometheus.Port), nil)
	if err != nil {
		fmt.Printf("error serving http: %v", err)
		return
	}
}
