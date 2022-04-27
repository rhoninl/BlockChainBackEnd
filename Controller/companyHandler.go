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
	if thisCompany.CompanyId == targetCompany.CompanyId {
		c.JSON(http.StatusBadRequest, gin.H{"message": "自己跟自己建交，你这个人可真有意思"})
		return
	}
	if !Model.CheckCompanyFriend(thisCompany.CompanyId, targetCompany.CompanyId) {
		c.JSON(http.StatusBadRequest, gin.H{"message": "都建交过了，还想咋建交，服了"})
		return
	}
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

func ReplyFriend(c *gin.Context) {
	companyId, _ := c.Get("companyId")
	var reply Utils.ReplyFriend
	c.Bind(&reply)
	reply.CompanyId = companyId.(int64)

	messageInfo, err := Model.GetMessageBasicInfo(reply.MessageId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "出错了，请联系工作人员"})
		return
	} else if reply.CompanyId != messageInfo.ToId {
		c.JSON(http.StatusBadRequest, gin.H{"message": "不是你的消息，瞎点nm"})
		return
	} else if messageInfo.IsReply == 1 {
		c.JSON(http.StatusBadRequest, gin.H{"message": "消息回过了还回，是不是有点大病?"})
		return
	}
	if reply.Ok {
		Model.PassReply(reply)
	} else {
		Model.SendMessageTo(0, "对方觉得你是个煞笔，所以拒绝了你的好友请求", messageInfo.ToId, messageInfo.FromId)
	}
	go Model.SetReply(reply.MessageId)
	c.JSON(http.StatusOK, nil)
}
