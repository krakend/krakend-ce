// +build linux

package main

import (
	"os"

	"github.com/devopsfaith/krakend/config"
	"github.com/devopsfaith/krakend/logging"
	"github.com/devopsfaith/krakend/plugin"
)

func loadPlugins(reg *plugin.Register, cfg config.ServiceConfig, logger logging.Logger) {
	if "" != os.Getenv("KRAKEND_ENABLE_PLUGINS") && cfg.Plugin != nil {
		logger.Info("Plugin experiment enabled!")
		pluginsLoaded, err := plugin.Load(*cfg.Plugin, reg)
		if err != nil {
			logger.Error(err.Error())
		}
		logger.Info("Total plugins loaded:", pluginsLoaded)
	}
}
