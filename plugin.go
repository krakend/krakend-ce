package krakend

import (
	"github.com/luraproject/lura/v2/logging"
	proxy "github.com/luraproject/lura/v2/proxy/plugin"
	client "github.com/luraproject/lura/v2/transport/http/client/plugin"
	server "github.com/luraproject/lura/v2/transport/http/server/plugin"
)

// LoadPlugins loads and registers the plugins so they can be used if enabled at the configuration
func LoadPlugins(folder, pattern string, logger logging.Logger) {
	logger.Debug("[SERVICE: Plugin Loader] Starting loading process")

	n, err := client.LoadWithLogger(
		folder,
		pattern,
		client.RegisterClient,
		logger,
	)
	logPluginLoaderErrors(logger, "[SERVICE: Executor Plugin]", n, err)

	n, err = server.LoadWithLogger(
		folder,
		pattern,
		server.RegisterHandler,
		logger,
	)
	logPluginLoaderErrors(logger, "[SERVICE: Handler Plugin]", n, err)

	n, err = proxy.LoadWithLogger(
		folder,
		pattern,
		proxy.RegisterModifier,
		logger,
	)
	logPluginLoaderErrors(logger, "[SERVICE: Modifier Plugin]", n, err)

	logger.Debug("[SERVICE: Plugin Loader] Loading process completed")
}

func logPluginLoaderErrors(logger logging.Logger, tag string, n int, err error) {
	if err != nil {
		if mErrs, ok := err.(pluginLoaderErr); ok {
			for _, err := range mErrs.Errs() {
				logger.Debug(tag, err.Error())
			}
		} else {
			logger.Debug(tag, err.Error())
		}
	}
	if n > 0 {
		logger.Info(tag, "Total plugins loaded:", n)
	}
}

type pluginLoader struct{}

func (d pluginLoader) Load(folder, pattern string, logger logging.Logger) {
	LoadPlugins(folder, pattern, logger)
}

type pluginLoaderErr interface {
	Errs() []error
}
