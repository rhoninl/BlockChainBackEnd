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

func GetStaff(c *gin.Context) {
	companyId, _ := c.Get("companyId")
	staffs, err := Model.GetStaff(companyId.(int64))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "服务器异常"})
		return
	}
	c.JSON(http.StatusOK, staffs)
}

func AddStaff(c *gin.Context) {
	companyId, _ := c.Get("companyId")
	var staffInfo Utils.Staff
	c.BindJSON(&staffInfo)
	id, err := Model.InsertStaff(staffInfo, companyId.(int64))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "服务器异常"})
		return
	} else if id == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"message": "该员工已存在"})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"staffId": id})
}

func DeleteStaff(c *gin.Context) {
	companyId, _ := c.Get("companyId")
	staffId := c.Query("id")
	iStaffId, err := strconv.ParseInt(staffId, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "请求错误"})
		return
	}
	if !Model.CheckStaffCompany(iStaffId, companyId.(int64)) {
		c.JSON(http.StatusBadRequest, gin.H{"message": "该员工不属于当前帐号"})
		return
	}
	if Model.DeleteStaff(iStaffId) != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "服务器异常"})
		return
	}
	c.JSON(888, nil)
}
func GetStaffInfo(c *gin.Context) {
	iCompanyId, _ := c.Get("companyId")
	companyId := iCompanyId.(int64)
	sStaffId := c.Query("staffId")
	staffId, err := strconv.ParseInt(sStaffId, 10, 64)
	if err != nil || staffId == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"message": "请求参数异常", "Query": staffId})
		return
	}
	if !Model.CheckStaffCompany(staffId, companyId) {
		c.JSON(http.StatusBadRequest, gin.H{"message": "该员工不属于当前帐号"})
		return
	}
	info, _, err := Model.GetStaffInfo(staffId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "服务器异常"})
		return
	}
	info.StaffId = staffId
	c.JSON(http.StatusOK, info)
}

func ChangeStaffInfo(c *gin.Context) {
	iCompanyId, _ := c.Get("companyId")
	companyId := iCompanyId.(int64)
	var info Utils.StaffInfo
	c.Bind(&info)
	if !Model.CheckStaffCompany(info.StaffId, companyId) {
		c.JSON(http.StatusBadRequest, gin.H{"message": "该员工不属于当前帐号"})
		return
	}
	oldInfo, addressId, err := Model.GetStaffInfo(info.StaffId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "服务器异常"})
		return
	}
	if reflect.DeepEqual(oldInfo, Info) {
		c.JSON(http.StatusNotModified, gin.H{"message": "信息没有改变"})
		return
	}
	if addressId == 0 || !reflect.DeepEqual(oldInfo.AddressInfo, info.AddressInfo) && !Model.UpdateStaffAddressInfo(info.AddressInfo, addressId, info.StaffId) {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "服务器异常"})
		return
	}
	oldInfo.AddressInfo = info.AddressInfo
	oldInfo.StaffId = info.StaffId
	oldInfo.JoinDate = ""
	if !reflect.DeepEqual(oldInfo, info) && !Model.UpdateStaffInfo(info) {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "服务器异常"})
		return
	}
	c.JSON(http.StatusCreated, nil)
}
