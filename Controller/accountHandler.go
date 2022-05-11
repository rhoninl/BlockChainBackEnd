package Controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"log"
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
	userInfo.Password, err = Utils.AesDecryptCBC(userInfo.Password)
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
		Secure:   true,
		HttpOnly: false,
		SameSite: 4,
	})
	c.JSON(http.StatusCreated, nil)
}

func LogOut(c *gin.Context) {
	companyId, _ := c.Get("companyId")
	http.SetCookie(c.Writer, &http.Cookie{
		Name:     "token",
		Value:    "Bye",
		Path:     "/",
		Domain:   "",
		MaxAge:   -1,
		Secure:   true,
		SameSite: 4,
	})
	go Model.UseClient().UnRegister(companyId.(int64))
	fmt.Println(companyId, "is LogOut")
	c.JSON(http.StatusOK, nil)
}

func Register(c *gin.Context) {
	var accountInfo Utils.RegisterInfo
	c.BindJSON(&accountInfo)
	password, err := Utils.AesDecryptCBC(accountInfo.Password)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "密码错误"})
		return
	}
	accountInfo.Password = password
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
	c.JSON(http.StatusCreated, nil)
}

func Info(c *gin.Context) {
	companyId, _ := c.Get("companyId")
	info, err := Model.Info(companyId.(int64))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "服务器出错了"})
		return
	}
	info.AddressId = 0
	c.JSON(http.StatusOK, info)
}

func GetAuth(c *gin.Context) {
	var info Utils.GetAuth
	err := c.Bind(&info)
	if err != nil || info.Email == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "邮箱无效"})
		return
	}
	if info.Tag != 0 && !Model.CheckEmailUnique(info.Email) {
		c.JSON(http.StatusBadRequest, gin.H{"message": "邮箱已使用"})
		log.Println(info)
		return
	}
	var codeInfo Utils.AuthCode
	codeInfo.Code = Utils.GenVerCode()
	codeInfo.ToEmail = info.Email
	err = Utils.AuthCodeRegister(codeInfo)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "验证码发送失败"})
		return
	}
	message := `<html><body>您的验证码为<h3>` + codeInfo.Code + `</h3><br/>验证码有效期为1小时，请在1小时内完成验证<br/>如果不是您本人操作，请忽略本条邮件</body></html>`
	if err := Utils.SendMessage(message, info.Email); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "服务器异常"})
		return
	}
	c.JSON(http.StatusCreated, nil)
}

func EditInfo(c *gin.Context) {
	companyId, _ := c.Get("companyId")
	var info Utils.Info
	c.Bind(&info)
	info.CompanyBasicInfo.CompanyId = companyId.(int64)
	info.CompanyInfo.CompanyId = info.CompanyBasicInfo.CompanyId
	try1 := Model.TryUpdateCompany(info.CompanyBasicInfo)
	try2 := Model.TryUpdateCompanyInfo(info.CompanyInfo)
	try3 := Model.TryUpdateAddress(info.AddressInfo, companyId.(int64))
	if try1 || try2 || try3 {
		c.JSON(http.StatusCreated, gin.H{"message": "修改成功"})
		return
	}
	c.JSON(http.StatusNotModified, gin.H{"message": "未做出修改"})
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
	c.JSON(http.StatusCreated, nil)
}

func ChangePassword(c *gin.Context) {
	var accountInfo Utils.RegisterInfo
	c.BindJSON(&accountInfo)
	companyId, _ := c.Get("companyId")
	accountInfo.Password, _ = Utils.AesDecryptCBC(accountInfo.Password)
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
	c.JSON(http.StatusCreated, nil)
}
