package Controller

import (
	"github.com/gin-gonic/gin"
	"math/rand"
	"net/http"
)

func GetIndex(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"chainHigh": rand.Intn(12345) >> 1,
		"chainNum":  rand.Intn(12345) >> 1,
		"nodeNum":   rand.Intn(12345) >> 1,
		"chainCode": rand.Intn(12345) >> 1,
	})
}
