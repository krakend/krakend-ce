// +build !linux

package main

import (
	"github.com/devopsfaith/krakend/config"
	"github.com/devopsfaith/krakend/logging"
	"github.com/devopsfaith/krakend/plugin"
)

func loadPlugins(_ *plugin.Register, _ config.ServiceConfig, _ logging.Logger) {}
