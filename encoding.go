package krakend

import (
	rss "github.com/devopsfaith/krakend-rss/v2"
	ginxml "github.com/devopsfaith/krakend-xml/gin"
	xml "github.com/devopsfaith/krakend-xml/v2"
	"github.com/luraproject/lura/router/gin"
)

// RegisterEncoders registers all the available encoders
func RegisterEncoders() {
	xml.Register()
	rss.Register()

	gin.RegisterRender(xml.Name, ginxml.Render)
}
