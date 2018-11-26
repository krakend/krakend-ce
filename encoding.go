package krakend

import (
	"github.com/devopsfaith/krakend-rss"
	"github.com/devopsfaith/krakend-xml"
)

// RegisterEncoders registers all the available encoders
func RegisterEncoders() {
	xml.Register()
	rss.Register()
}
