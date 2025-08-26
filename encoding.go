package krakend

import (
	rss "github.com/krakend/krakend-rss/v2"
	xml "github.com/krakend/krakend-xml/v2"
	ginxml "github.com/krakend/krakend-xml/v2/gin"
	"github.com/luraproject/lura/v2/router/gin"
)

// RegisterEncoders registers all the available encoders
func RegisterEncoders() {
	xml.Register()
	rss.Register()

	gin.RegisterRender(xml.Name, ginxml.Render)
}
