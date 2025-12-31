package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/go-viper/mapstructure/v2"
)

var pluginName = "hog-static-content"

var HandlerRegisterer = registerer(pluginName)

type registerer string

func (r registerer) RegisterHandlers(f func(
	name string,
	handler func(context.Context, map[string]interface{}, http.Handler) (http.Handler, error),
)) {
	f(string(r), r.registerHandlers)
}

// The static content plugin will look for this configuration:
/*
	   "extra_config": {
	       "plugin/http-server": {
	           "name":["hog-static-content"],
	           "hog-static-content": {
	               "static": [{
						"path-prefix": "/*",
						"service-host": "http://web-example"
						"keep-unsafe-headers": false
					}],
					"service-gateway": {
						"path-prefix": ["/api/*"],
					}
	           }
	       }
	   }
*/
func (r registerer) registerHandlers(_ context.Context, extra map[string]interface{}, h http.Handler) (http.Handler, error) {
	logger.Debug(fmt.Sprintf("Loading static-content plugin config"))
	config, err := loadPluginConfig(extra)
	if err != nil {
		return nil, errors.Join(errors.New("failed to load plugin config"), err)
	}

	logger.Debug(fmt.Sprintf("Registering static content routes"))
	for _, s := range config.Static {
		logger.Debug(fmt.Sprintf("The plugin is now servig static content on the path %s", s.PathPrefix))
	}

	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		handleStaticContent(config, writer, request, h)
	}), nil
}

func main() {

}

type PluginConfig struct {
	Static         []StaticConfig       `mapstructure:"static"`
	ServiceGateway ServiceGatewayConfig `mapstructure:"service-gateway"`
}

type StaticConfig struct {
	PathPrefix        string `mapstructure:"path-prefix"`
	ServiceHost       string `mapstructure:"service-host"`
	KeepUnsafeHeaders bool   `mapstructure:"keep-unsafe-headers"`
}

type ServiceGatewayConfig struct {
	PathPrefix []string `mapstructure:"path-prefix"`
}

func loadPluginConfig(cfg map[string]interface{}) (PluginConfig, error) {
	var pc PluginConfig
	err := mapstructure.Decode(cfg[pluginName], &pc)
	if err != nil {
		return pc, fmt.Errorf("failed to decode config: %w", err)
	}
	return pc, nil
}

// This logger is replaced by the RegisterLogger method to load the one from KrakenD
var logger Logger = noopLogger{}

func (registerer) RegisterLogger(v interface{}) {
	l, ok := v.(Logger)
	if !ok {
		return
	}
	logger = l
	logger.Debug(fmt.Sprintf("[PLUGIN: %s] Logger loaded", HandlerRegisterer))
}

type Logger interface {
	Debug(v ...interface{})
	Info(v ...interface{})
	Warning(v ...interface{})
	Error(v ...interface{})
	Critical(v ...interface{})
	Fatal(v ...interface{})
}

// Empty logger implementation
type noopLogger struct{}

func (n noopLogger) Debug(_ ...interface{})    {}
func (n noopLogger) Info(_ ...interface{})     {}
func (n noopLogger) Warning(_ ...interface{})  {}
func (n noopLogger) Error(_ ...interface{})    {}
func (n noopLogger) Critical(_ ...interface{}) {}
func (n noopLogger) Fatal(_ ...interface{})    {}
