package krakend

import (
	"context"
	"fmt"

	consul "github.com/devopsfaith/krakend-consul"
	"github.com/devopsfaith/krakend/config"
	"github.com/devopsfaith/krakend/logging"
	"github.com/devopsfaith/krakend/sd/dnssrv"
)

// RegisterSubscriberFactories registers all the available sd adaptors
func RegisterSubscriberFactories(ctx context.Context, cfg config.ServiceConfig, logger logging.Logger) func(n string, p int) {
	// register the dns service discovery
	dnssrv.Register()

	return func(name string, port int) {
		if err := consul.Register(ctx, cfg.ExtraConfig, port, name, logger); err != nil {
			logger.Error(fmt.Sprintf("Couldn't register %s:%d in consul: %s", name, port, err.Error()))
		}
	}
}

type registerSubscriberFactories struct{}

func (d registerSubscriberFactories) Register(ctx context.Context, cfg config.ServiceConfig, logger logging.Logger) func(n string, p int) {
	return RegisterSubscriberFactories(ctx, cfg, logger)
}
