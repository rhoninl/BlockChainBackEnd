package Controller

import (
	"github.com/gin-gonic/gin"
	"main/Model"
	"net/http"
)

func GetAllOrder(c *gin.Context) {
	companyId, _ := c.Get("companyId")
	orderInfo, err := Model.GetAllOrder(companyId.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "服务器异常"})
		return
	}
	c.JSON(http.StatusOK, orderInfo)
}
