package krakend

import (
	"context"
	"fmt"

	consul "github.com/devopsfaith/krakend-consul"
	"github.com/devopsfaith/krakend-etcd"
	"github.com/devopsfaith/krakend/config"
	"github.com/devopsfaith/krakend/logging"
	"github.com/devopsfaith/krakend/sd"
	"github.com/devopsfaith/krakend/sd/dnssrv"
)

// RegisterSubscriberFactories registers all the available sd adaptors
func RegisterSubscriberFactories(ctx context.Context, cfg config.ServiceConfig, logger logging.Logger) func(n string, p int) {
	// setup the etcd client if necessary
	etcdClient, err := etcd.New(ctx, cfg.ExtraConfig)
	if err != nil {
		logger.Warning("building the etcd client:", err.Error())
	}
	sd.RegisterSubscriberFactory("etcd", etcd.SubscriberFactory(ctx, etcdClient))

	// register the dns service discovery
	dnssrv.Register()

	return func(name string, port int) {
		if err := consul.Register(ctx, cfg.ExtraConfig, port, name, logger); err != nil {
			logger.Error(fmt.Sprintf("Couldn't register %s:%d in consul: %s", name, port, err.Error()))
		}

		// TODO: add the call to the etcd service register
	}
}
