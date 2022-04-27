package Controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
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
	if err != nil || bcrypt.CompareHashAndPassword([]byte(currentInfo.Password), []byte(userInfo.Password)) != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "密码不正确"})
		return
	}
	token := Utils.CreateToken(currentInfo.CompanyId)
	//c.SetCookie("token", token, Utils.MAXAGE, "/", "", true, false)
	http.SetCookie(c.Writer, &http.Cookie{
		Name:     "token",
		Value:    token,
		Path:     "/",
		Domain:   "",
		MaxAge:   Utils.MAXAGE,
		Secure:   false,
		HttpOnly: false,
		SameSite: 4,
	})
	c.JSON(http.StatusOK, nil)
}

func LogOut(c *gin.Context) {
	companyId, _ := c.Get("companyId")

	go Model.UseClient().UnRegister(companyId.(int64))
	fmt.Println(companyId, "is LogOut")
}

func Register(c *gin.Context) {
	var accountInfo Utils.RegisterInfo
	c.BindJSON(&accountInfo)
	//password, err := Utils.ParsePassword(accountInfo.Password)
	//if err != nil {
	//	c.JSON(http.StatusBadRequest, gin.H{"message": "密码错误"})
	//	return
	//}
	//accountInfo.Password = password
	//if !Model.CheckEmailUnique(accountInfo.ToEmail) {
	//	c.JSON(http.StatusBadRequest, gin.H{"message": "邮箱已被使用"})
	//	return
	//}
	//if !Model.CheckAccountUnique(accountInfo.Account.Account) {
	//	c.JSON(http.StatusBadRequest, gin.H{"message": "帐号已存在"})
	//	return
	//}
	//ok, err := Utils.AuthCodeCheck(accountInfo.AuthCode)
	//if err != nil {
	//	c.JSON(http.StatusInternalServerError, gin.H{"message": "服务器异常"})
	//	return
	//} else if !ok {
	//	c.JSON(http.StatusBadRequest, gin.H{"message": "验证码错误"})
	//	return
	//}
	ok, err := Model.RegisterInfo(accountInfo)
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
	var info Utils.GetAuth
	c.Bind(&info)
	if info.Email == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "邮箱无效"})
		return
	}
	if info.Tag != "" && !Model.CheckEmailUnique(info.Email) {
		c.JSON(http.StatusBadRequest, gin.H{"message": "邮箱已使用"})
		return
	}
	var codeInfo Utils.AuthCode
	codeInfo.Code = Utils.GenVerCode()
	message := `<html><body style:"background:'https://gimg2.baidu.com/image_search/src=http%3A%2F%2Fbkimg.cdn.bcebos.com%2Fpic%2Fb3fb43166d224f4a20a4652df4a687529822720e7bc9&refer=http%3A%2F%2Fbkimg.cdn.bcebos.com&app=2002&size=f9999,10000&q=a80&n=0&g=0n&fmt=auto?sec=1653183161&t=51bbd5bef8fe26d143d2fd2a4a73abf3'"><a>您的验证码为</a><h3>` + codeInfo.Code + `</h3><a><br/>验证码有效期为1小时，请在1小时内完成验证<br/>如果不是您本人操作，请忽略本条邮件</a></body></html>`
	if err := Utils.SendMessage(message, info.Email); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "服务器异常"})
		return
	}
	codeInfo.ToEmail = info.Email
	Utils.AuthCodeRegister(codeInfo)
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

func ForgetPassword(c *gin.Context) {
	var form Utils.ForgetPasswordForm
	c.Bind(&form)
	if !Model.CheckEmail(form.Account, form.Email) {
		c.JSON(http.StatusBadRequest, gin.H{"message": "帐号或者邮箱错误"})
		return
	}
	newPassword := Utils.GeneratePassWord()
	if !Model.ChangePassword(form.Account, newPassword) {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "服务器异常"})
		return
	}
	message := `<html><body>您的新密码为<h3>` + newPassword + `</h3><br/>登陆后请及时修改<br/></body></html>`
	Utils.SendMessage(message, form.Email)
	c.JSON(http.StatusOK, nil)
}

func ChangePassword(c *gin.Context) {
	var accountInfo Utils.RegisterInfo
	c.BindJSON(&accountInfo)
	companyId, _ := c.Get("companyId")
	accountInfo.Password, _ = Utils.ParsePassword(accountInfo.Password)
	accountInfo.Account.CompanyId = companyId.(int64)
	if !Model.CheckEmail(accountInfo.Account.CompanyId, accountInfo.ToEmail) {
		c.JSON(http.StatusBadRequest, gin.H{"message": "帐号或者邮箱错误"})
		return
	}
	ok, err := Utils.AuthCodeCheck(accountInfo.AuthCode)
	if !ok || err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "请求异常"})
		return
	}
	if !Model.ChangePassword(accountInfo.Account.CompanyId, accountInfo.Account.Password) {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "服务器异常"})
		return
	}
	c.JSON(http.StatusOK, nil)
}
