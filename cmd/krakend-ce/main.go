// Krakend-ce sets up a complete KrakenD API Gateway ready to serve

package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	krakend "github.com/devopsfaith/krakend-ce"
	"github.com/devopsfaith/krakend-cobra"
	flexibleconfig "github.com/devopsfaith/krakend-flexibleconfig"
	"github.com/devopsfaith/krakend-viper"
	"github.com/devopsfaith/krakend/config"
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
