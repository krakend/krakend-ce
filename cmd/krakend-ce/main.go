// Krakend-ce sets up a complete KrakenD API Gateway ready to serve

package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	krakend "github.com/krakendio/krakend-ce/v2"
	cmd "github.com/krakendio/krakend-cobra/v2"
	flexibleconfig "github.com/krakendio/krakend-flexibleconfig/v2"
	viper "github.com/krakendio/krakend-viper/v2"
	"github.com/luraproject/lura/v2/config"
)

const (
	fcPartials  = "FC_PARTIALS"
	fcTemplates = "FC_TEMPLATES"
	fcSettings  = "FC_SETTINGS"
	fcPath      = "FC_OUT"
	fcEnable    = "FC_ENABLE"
)

func main() {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		select {
		case sig := <-sigs:
			log.Println("Signal intercepted:", sig)
			cancel()
		case <-ctx.Done():
		}
	}()

	krakend.RegisterEncoders()

	for key, alias := range aliases {
		config.ExtraConfigAlias[alias] = key
	}

	var cfg config.Parser
	cfg = viper.New()
	if os.Getenv(fcEnable) != "" {
		cfg = flexibleconfig.NewTemplateParser(flexibleconfig.Config{
			Parser:    cfg,
			Partials:  os.Getenv(fcPartials),
			Settings:  os.Getenv(fcSettings),
			Path:      os.Getenv(fcPath),
			Templates: os.Getenv(fcTemplates),
		})
	}

	cmd.Execute(cfg, krakend.NewExecutor(ctx))
}

var aliases = map[string]string{
	"github_com/devopsfaith/krakend/transport/http/server/handler":  "plugin/http-server",
	"github.com/devopsfaith/krakend/transport/http/client/executor": "plugin/http-client",
	"github.com/devopsfaith/krakend/proxy/plugin":                   "plugin/req-resp-modifier",
	"github.com/devopsfaith/krakend/proxy":                          "proxy",
	"github_com/luraproject/lura/router/gin":                        "router",

	"github.com/devopsfaith/krakend-ratelimit/juju/router":    "qos/ratelimit/router",
	"github.com/devopsfaith/krakend-ratelimit/juju/proxy":     "qos/ratelimit/proxy",
	"github.com/devopsfaith/krakend-httpcache":                "qos/http-cache",
	"github.com/devopsfaith/krakend-circuitbreaker/gobreaker": "qos/circuit-breaker",

	"github.com/devopsfaith/krakend-oauth2-clientcredentials": "auth/client-credentials",
	"github.com/devopsfaith/krakend-jose/validator":           "auth/validator",
	"github.com/devopsfaith/krakend-jose/signer":              "auth/signer",
	"github_com/devopsfaith/bloomfilter":                      "auth/revoker",

	"github_com/devopsfaith/krakend-botdetector": "security/bot-detector",
	"github_com/devopsfaith/krakend-httpsecure":  "security/http",
	"github_com/devopsfaith/krakend-cors":        "security/cors",

	"github.com/devopsfaith/krakend-cel":        "validation/cel",
	"github.com/devopsfaith/krakend-jsonschema": "validation/json-schema",

	"github.com/devopsfaith/krakend-amqp/agent": "async/amqp",

	"github.com/devopsfaith/krakend-amqp/consume":                  "backend/amqp/consumer",
	"github.com/devopsfaith/krakend-amqp/produce":                  "backend/amqp/producer",
	"github.com/devopsfaith/krakend-lambda":                        "backend/lambda",
	"github.com/devopsfaith/krakend-pubsub/publisher":              "backend/pubsub/publisher",
	"github.com/devopsfaith/krakend-pubsub/subscriber":             "backend/pubsub/subscriber",
	"github.com/devopsfaith/krakend/transport/http/client/graphql": "backend/graphql",
	"github.com/devopsfaith/krakend/http":                          "backend/http",

	"github_com/devopsfaith/krakend-gelf":       "telemetry/gelf",
	"github_com/devopsfaith/krakend-gologging":  "telemetry/logging",
	"github_com/devopsfaith/krakend-logstash":   "telemetry/logstash",
	"github_com/devopsfaith/krakend-metrics":    "telemetry/metrics",
	"github_com/letgoapp/krakend-influx":        "telemetry/influx",
	"github_com/devopsfaith/krakend-influx":     "telemetry/influx",
	"github_com/Jozefiel/krakend-influx2":       "telemetry/influx2",
	"github_com/devopsfaith/krakend-opencensus": "telemetry/opencensus",

	"github.com/devopsfaith/krakend-lua/router":        "modifier/lua-endpoint",
	"github.com/devopsfaith/krakend-lua/proxy":         "modifier/lua-proxy",
	"github.com/devopsfaith/krakend-lua/proxy/backend": "modifier/lua-backend",
	"github.com/devopsfaith/krakend-martian":           "modifier/martian",
}
