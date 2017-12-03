package main

import (
	"context"

	"github.com/devopsfaith/krakend-etcd"
	"github.com/devopsfaith/krakend/config"
	"github.com/devopsfaith/krakend/logging"
	"github.com/devopsfaith/krakend/sd"
	"github.com/devopsfaith/krakend/sd/dnssrv"
)

// RegisterSubscriberFactories registers all the available sd adaptors
func RegisterSubscriberFactories(ctx context.Context, cfg config.ServiceConfig, logger logging.Logger) {
	// setup the etcd client if necessary
	etcdClient, err := etcd.New(ctx, cfg.ExtraConfig)
	if err != nil {
		logger.Warning("building the etcd client:", err.Error())
	}
	sd.RegisterSubscriberFactory("etcd", etcd.SubscriberFactory(ctx, etcdClient))

	// register the dns service discovery
	dnssrv.Register()
}
