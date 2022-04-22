package Controller

import (
	"github.com/gin-gonic/gin"
	"main/Model"
	"net/http"
)

func GetAllMessage(c *gin.Context) {
	companyId, _ := c.Get("companyId")
	message, err := Model.GetAllMessage(companyId.(int64))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "服务器异常"})
		return
	}
	c.JSON(http.StatusOK, message)
}
