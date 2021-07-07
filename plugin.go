package krakend

import (
	"github.com/luraproject/lura/logging"
	proxy "github.com/luraproject/lura/proxy/plugin"
	client "github.com/luraproject/lura/transport/http/client/plugin"
	server "github.com/luraproject/lura/transport/http/server/plugin"
)

// LoadPlugins loads and registers the plugins so they can be used if enabled at the configuration
func LoadPlugins(folder, pattern string, logger logging.Logger) {
	n, err := client.LoadWithLogger(
		folder,
		pattern,
		client.RegisterClient,
		logger,
	)
	if err != nil {
		logger.Warning("loading plugins:", err)
	}
	logger.Info("total http executor plugins loaded:", n)

	n, err = server.LoadWithLogger(
		folder,
		pattern,
		server.RegisterHandler,
		logger,
	)
	if err != nil {
		logger.Warning("loading plugins:", err)
	}
	logger.Info("total http handler plugins loaded:", n)

	n, err = proxy.LoadModifiersWithLogger(
		folder,
		pattern,
		proxy.RegisterModifier,
		logger,
	)
	if err != nil {
		logger.Warning("loading plugins:", err)
	}
	logger.Info("total modifier plugins loaded:", n)
}

type pluginLoader struct{}

func (d pluginLoader) Load(folder, pattern string, logger logging.Logger) {
	LoadPlugins(folder, pattern, logger)
}
