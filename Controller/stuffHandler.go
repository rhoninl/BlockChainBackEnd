package Controller

import (
	"github.com/gin-gonic/gin"
	"main/Model"
	"main/Utils"
	"net/http"
	"strconv"
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
	c.BindJSON(&stuffInfo)
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

func DeleteStuff(c *gin.Context) {
	companyId, _ := c.Get("companyId")
	stuffId := c.Query("id")
	iStuffId, err := strconv.ParseInt(stuffId, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "请求错误"})
		return
	}
	if !Model.CheckStuffCompany(iStuffId, companyId.(int64)) {
		c.JSON(http.StatusBadRequest, gin.H{"message": "不是你的员工，删nm"})
		return
	}
	if Model.DeleteStuff(iStuffId) != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "服务器异常"})
		return
	}
	c.JSON(888, nil)
}
