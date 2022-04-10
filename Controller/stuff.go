package Controller

import (
	"github.com/gin-gonic/gin"
	"main/Model"
	"net/http"
)

func GetStuff(c *gin.Context) {
	companyId, _ := c.Get("companyId")
	stuffs, err := Model.GetStuff(companyId.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "服务器异常"})
		return
	}
	c.JSON(http.StatusOK, stuffs)
}
