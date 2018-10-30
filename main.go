package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/devopsfaith/krakend-cobra"
	flexibleconfig "github.com/devopsfaith/krakend-flexibleconfig"
	"github.com/devopsfaith/krakend-viper"
)

const (
	fcPartials = "FC_PARTIALS"
	fcSettings = "FC_SETTINGS"
	fcPath     = "FC_OUT"
	fcEnable   = "FC_ENABLE"
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

	RegisterEncoders()

	cfg := viper.New()
	if os.Getenv(fcEnable) != "" {
		cfg = flexibleconfig.NewTemplateParser(flexibleconfig.Config{
			Parser:   cfg,
			Partials: os.Getenv(fcPartials),
			Settings: os.Getenv(fcSettings),
			Path:     os.Getenv(fcPath),
		})
	}

	cmd.Execute(cfg, NewExecutor(ctx))
}
