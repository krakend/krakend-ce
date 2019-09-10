package krakend

import (
	"net/http"
	"github.com/gin-gonic/gin"
	"github.com/devopsfaith/krakend/proxy"
)

func jsonErrorRender(c *gin.Context, response *proxy.Response) {
	if response == nil {
		c.JSON(http.StatusOK, gin.H{})
	}
	status := http.StatusOK
	if response.Metadata.StatusCode > 0 {
		status = response.Metadata.StatusCode
	}
	c.JSON(status, response.Data)
}
