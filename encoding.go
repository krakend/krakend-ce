package krakend

import (
	rss "github.com/krakend/krakend-rss/v3"
	xml "github.com/krakend/krakend-xml/v3"
	ginxml "github.com/krakend/krakend-xml/v3/gin"
	"github.com/luraproject/lura/v3/router/gin"
)

// RegisterEncoders registers all the available encoders
func RegisterEncoders() {
	xml.Register()
	rss.Register()

	gin.RegisterRender(xml.Name, ginxml.Render)
}
