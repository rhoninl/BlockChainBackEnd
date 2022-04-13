package Controller

import (
	"github.com/gin-gonic/gin"
	"main/Model"
	"main/Utils"
	"net/http"
)

func GetStuff(c *gin.Context) {
	companyId, _ := c.Get("companyId")
	stuffs, err := Model.GetStuff(companyId.(int64))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "服务器异常"})
		return
	}
	c.JSON(http.StatusOK, stuffs)
}

func AddStuff(c *gin.Context) {
	companyId, _ := c.Get("companyId")
	var stuffInfo Utils.Stuff
	c.Bind(&stuffInfo)
	id, err := Model.InsertStuff(stuffInfo, companyId.(int64))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "服务器异常"})
		return
	} else if id == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"message": "该员工已存在"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"stuffId": id})
}
