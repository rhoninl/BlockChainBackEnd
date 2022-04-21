package Controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"main/Model"
	"main/Utils"
	"net/http"
)

func Login(c *gin.Context) {
	var userInfo Utils.Account
	c.BindJSON(&userInfo)
	if userInfo.Account == "" || userInfo.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "帐号或者密码不能为空"})
		return
	}
	currentInfo, exists, err := Model.Login(userInfo)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "服务器出错了"})
		return
	}
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"message": "帐号不存在"})
		return
	}
	userInfo.Password, err = Utils.ParsePassword(userInfo.Password)
	if err != nil || currentInfo.Password != userInfo.Password {
		c.JSON(http.StatusBadRequest, gin.H{"message": "密码不正确"})
		return
	}
	token := Utils.CreateToken(currentInfo.CompanyId)
	c.SetCookie("token", token, Utils.MAXAGE, "/", "", false, true)
	c.JSON(http.StatusOK, nil)
}

func LogOut(c *gin.Context) {
	companyId, _ := c.Get("companyId")
	c.SetCookie("token", "", -1, "/", "", false, true)
	UseClient().UnRegister(companyId.(int64))
	fmt.Println(companyId, "is LogOut")
}

func Register(c *gin.Context) {
	var accountInfo Utils.RegisterInfo
	c.BindJSON(&accountInfo)
	if !Model.CheckEmailUnique(accountInfo.ToEmail) {
		c.JSON(http.StatusBadRequest, gin.H{"message": "邮箱已被使用"})
		return
	}
	if !Model.CheckAccountUnique(accountInfo.Account.Account) {
		c.JSON(http.StatusBadRequest, gin.H{"message": "帐号已存在"})
		return
	}
	ok, err := Utils.AuthCodeCheck(accountInfo.AuthCode)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "服务器异常"})
		return
	} else if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"message": "验证码错误"})
		return
	}
	ok, err = Model.RegisterInfo(accountInfo)
	if err != nil || !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "服务器异常"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "ok"})
}

func Info(c *gin.Context) {
	companyId, _ := c.Get("companyId")
	info, exists, err := Model.Info(companyId.(int64))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "服务器出错了"})
		return
	}
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"message": "未找到相关信息"})
		return
	}
	c.JSON(http.StatusOK, info)
}

func GetAuth(c *gin.Context) {
	var info Utils.AuthCode
	info.ToEmail = c.Query("email")
	if info.ToEmail == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "邮箱无效"})
		return
	}
	if !Model.CheckEmailUnique(info.ToEmail) {
		c.JSON(http.StatusBadRequest, gin.H{"message": "邮箱已使用"})
		return
	}
	info.Code = Utils.GenVerCode()
	if err := Utils.SendCode(info); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "服务器异常"})
		return
	}
	c.JSON(http.StatusOK, nil)
}

func EditInfo(c *gin.Context) {
	companyId, _ := c.Get("companyId")
	var companyInfo Utils.CompanyInfo
	c.Bind(&companyInfo)
	companyInfo.CompanyId = companyId.(int64)
	try1 := Model.TryUpdateCompany(companyInfo.CompanyBasicInfo)
	try2 := Model.TryUpdateCompanyInfo(companyInfo)
	try3 := Model.TryUpdateAddress(companyInfo.AddressInfo, companyInfo.CompanyId)
	if try1 || try2 || try3 {
		c.JSON(http.StatusOK, gin.H{"message": "修改成功"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "未做出修改"})
}
