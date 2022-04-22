package Controller

import (
	"github.com/gin-gonic/gin"
	"main/Model"
	"main/Utils"
	"net/http"
	"strconv"
)

func GetJointVenture(c *gin.Context) {
	companyId, _ := c.Get("companyId")
	companyList, err := Model.GetJointVenture(companyId.(int64))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "数据库异常"})
		return
	}
	c.JSON(http.StatusOK, companyList)
}

func MakeFriend(c *gin.Context) {
	var thisCompany, targetCompany Utils.CompanyBasicInfo
	fromCompanyId, _ := c.Get("companyId")
	c.Bind(&targetCompany)
	thisCompany.CompanyId = fromCompanyId.(int64)
	thisCompany.CompanyName, thisCompany.CompanyType = Model.GetCompanyBasicInfo(thisCompany.CompanyId)
	targetCompany.CompanyName, _ = Model.GetCompanyBasicInfo(targetCompany.CompanyId)
	if targetCompany.CompanyName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "公司不存在"})
		return
	}
	msg := thisCompany.CompanyName + `(` + strconv.FormatInt(thisCompany.CompanyId, 10) + `)想和你建交`
	if !Model.SendMessageTo(1, msg, targetCompany.CompanyId, thisCompany.CompanyId) {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "服务器异常"})
		return
	}
	c.JSON(http.StatusOK, nil)
}
