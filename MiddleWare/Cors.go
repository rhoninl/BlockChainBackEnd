package MiddleWare

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func MyCors() gin.HandlerFunc {
	return func(c *gin.Context) {
		config := cors.DefaultConfig()
		config.AllowAllOrigins = false
		config.AllowCredentials = true
		config.AllowOrigins = []string{c.Request.Header.Get("Origin")}
		cors.New(config)
	}
}
