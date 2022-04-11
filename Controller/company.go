package Controller

import (
	"github.com/gin-gonic/gin"
	"main/Model"
	"net/http"
)

func GetJointVenture(c *gin.Context) {
	companyId, _ := c.Get("companyId")
	companyList, err := Model.GetJointVenture(companyId.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "数据库异常"})
		return
	}
	c.JSON(http.StatusOK, companyList)
}
