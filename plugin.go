package krakend

import (
	"github.com/devopsfaith/krakend/logging"
	client "github.com/devopsfaith/krakend/transport/http/client/plugin"
	server "github.com/devopsfaith/krakend/transport/http/server/plugin"
)

// LoadPlugins loads and registers the plugins so they can be used if enabled at the configuration
func LoadPlugins(folder, pattern string, logger logging.Logger) {
	n, err := client.Load(
		folder,
		pattern,
		client.RegisterClient,
	)
	if err != nil {
		logger.Warning("loading plugins:", err)
	}
	logger.Info("total http executor plugins loaded:", n)

	n, err = server.Load(
		folder,
		pattern,
		server.RegisterHandler,
	)
	if err != nil {
		logger.Warning("loading plugins:", err)
	}
	logger.Info("total http handler plugins loaded:", n)
}
