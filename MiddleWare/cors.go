package MiddleWare

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func Co() gin.HandlerFunc {
	return func(c *gin.Context) {
		config := cors.DefaultConfig()
		config.AllowAllOrigins = false
		config.AllowOrigins = []string{c.Request.Header.Get("Origin")}
		cors.New(config)
	}
}
