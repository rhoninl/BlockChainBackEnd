package Controller

import (
	"github.com/gin-gonic/gin"
	"main/Model"
	"net/http"
)

func GetStuff(c *gin.Context) {
	companyId, exists := c.Get("companyId")
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"message": "未登录"})
		return
	}
	stuffs, err := Model.GetStuff(companyId.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "服务器异常"})
		return
	}
	c.JSON(http.StatusOK, stuffs)
}
