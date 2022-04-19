package Controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetIndex(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"chainHigh": 1,
		"chainNum":  2,
		"nodeNum":   3,
		"chainCode": 4,
	})
}
