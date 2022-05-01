package Controller

import "C"
import (
	"github.com/gin-gonic/gin"
	"main/Model"
	"main/Utils"
	"net/http"
	"reflect"
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
	c.JSON(http.StatusCreated, gin.H{"stuffId": id})
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
		c.JSON(http.StatusBadRequest, gin.H{"message": "该员工不属于当前帐号"})
		return
	}
	if Model.DeleteStuff(iStuffId) != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "服务器异常"})
		return
	}
	c.JSON(888, nil)
}
func GetStuffInfo(c *gin.Context) {
	iCompanyId, _ := c.Get("companyId")
	companyId := iCompanyId.(int64)
	sStuffId := c.Query("stuffId")
	stuffId, err := strconv.ParseInt(sStuffId, 10, 64)
	if err != nil || stuffId == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"message": "请求参数异常", "Query": stuffId})
		return
	}
	if !Model.CheckStuffCompany(stuffId, companyId) {
		c.JSON(http.StatusBadRequest, gin.H{"message": "该员工不属于当前帐号"})
		return
	}
	info, _, err := Model.GetStuffInfo(stuffId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "服务器异常"})
		return
	}
	info.StuffId = stuffId
	c.JSON(http.StatusOK, info)
}

func ChangeStuffInfo(c *gin.Context) {
	iCompanyId, _ := c.Get("companyId")
	companyId := iCompanyId.(int64)
	var info Utils.StuffInfo
	c.Bind(&info)
	if !Model.CheckStuffCompany(info.StuffId, companyId) {
		c.JSON(http.StatusBadRequest, gin.H{"message": "该员工不属于当前帐号"})
		return
	}
	oldInfo, addressId, err := Model.GetStuffInfo(info.StuffId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "服务器异常"})
		return
	}
	if reflect.DeepEqual(oldInfo, Info) {
		c.JSON(http.StatusNotModified, gin.H{"message": "信息没有改变"})
		return
	}
	if addressId == 0 || !reflect.DeepEqual(oldInfo.AddressInfo, info.AddressInfo) && !Model.UpdateStuffAddressInfo(info.AddressInfo, addressId, info.StuffId) {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "服务器异常"})
		return
	}
	oldInfo.AddressInfo = info.AddressInfo
	oldInfo.StuffId = info.StuffId
	oldInfo.JoinDate = ""
	if !reflect.DeepEqual(oldInfo, info) && !Model.UpdateStuffInfo(info) {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "服务器异常"})
		return
	}
	c.JSON(http.StatusCreated, nil)
}
