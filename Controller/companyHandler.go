package Controller

import (
	"github.com/gin-gonic/gin"
	"log"
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
		c.JSON(http.StatusBadRequest, gin.H{"message": "无法与自己建交"})
		return
	}
	if !Model.CheckCompanyFriend(thisCompany.CompanyId, targetCompany.CompanyId) {
		c.JSON(http.StatusBadRequest, gin.H{"message": "已与该公司建交"})
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
	c.JSON(http.StatusCreated, nil)
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
		c.JSON(http.StatusBadRequest, gin.H{"message": "该消息的接收帐号与当前帐号不匹配"})
		return
	} else if messageInfo.IsReply == 1 {
		c.JSON(http.StatusBadRequest, gin.H{"message": "该消息已回复"})
		return
	}
	if reply.Ok {
		Model.SendMessageTo(0, "对方通过了你的请求", messageInfo.FromId, 0)
		Model.PassReply(reply)
	} else {
		Model.SendMessageTo(0, "很抱歉地通知您，对方拒绝了你的好友请求", messageInfo.FromId, 0)
	}
	go Model.SetReply(reply.MessageId)
	c.JSON(http.StatusCreated, nil)
}

func DeleteFriend(c *gin.Context) {
	thisCompanyId, _ := c.Get("companyId")
	var info Utils.CompanyBasicInfo
	c.Bind(&info)
	if thisCompanyId.(int64) == info.CompanyId {
		c.JSON(http.StatusBadRequest, gin.H{"message": "跟自己绝交"})
		return
	}
	if !Model.CheckCompanyFriend(thisCompanyId.(int64), info.CompanyId) {
		c.JSON(http.StatusBadRequest, gin.H{"message": "删除的公司尚未建交"})
		return
	}
	if Model.DeleteCompanyFriend(thisCompanyId.(int64), info.CompanyId) != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "服务器异常"})
		return
	}
	info, _ = Model.CompanyBasicInfo(info.CompanyId)
	Model.SendMessageTo(0, "很抱歉通知您，您的好友"+info.CompanyName+"( id: "+strconv.FormatInt(info.CompanyId, 10)+" )把你给删了", info.CompanyId, 0)
	c.JSON(http.StatusCreated, nil)
}

func GetFriendsInfo(c *gin.Context) {
	thisCompanyId, _ := c.Get("companyId")
	targetCompanyId, err := strconv.ParseInt(c.Query("companyId"), 10, 64)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"message": "请求异常"})
		return
	}
	if !Model.CheckCompanyFriend(thisCompanyId.(int64), targetCompanyId) {
		c.JSON(http.StatusBadRequest, gin.H{"message": "当前公司与你尚未建交"})
		return
	}
	info, exists, err := Model.Info(targetCompanyId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "服务器异常"})
		return
	} else if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"message": "目标企业不存在"})
		return
	}
	c.JSON(http.StatusOK, info)
}

func SendMessageToFriends(c *gin.Context) {
	companyId, _ := c.Get("companyId")
	var info Utils.MessageStruct
	c.Bind(&info)
	if !Model.CheckCompanyFriend(companyId.(int64), info.ToId) {
		c.JSON(http.StatusBadRequest, gin.H{"message": "当前公司与你尚未建交"})
		return
	}
	if Model.SendMessageTo(2, info.Context, info.ToId, companyId.(int64)) {
		c.JSON(http.StatusCreated, nil)
	} else {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "消息未发送"})
	}
}
