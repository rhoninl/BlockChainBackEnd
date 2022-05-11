package Controller

import (
	"github.com/gin-gonic/gin"
	"main/Model"
	"net/http"
	"strconv"
)

func GetMessage(c *gin.Context) {
	companyId, _ := c.Get("companyId")
	messageId := c.Query("messageId")
	mid, err := strconv.ParseInt(messageId, 10, 64)
	if messageId == "" || err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "请求异常"})
		return
	}
	message, err := Model.GetMessage(companyId.(int64), mid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "服务器异常"})
		return
	}
	c.JSON(http.StatusOK, message)
}

func GetMessageInfo(c *gin.Context) {
	messageId := c.Query("id")
	if messageId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "信息不存在"})
		return
	}
	companyId, _ := c.Get("companyId")
	fMessageId, err := strconv.ParseInt(messageId, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "信息错误"})
		return
	}
	if !Model.CheckMessageAuth(fMessageId, companyId.(int64)) {
		c.JSON(http.StatusBadRequest, gin.H{"message": "信息的接受帐号与当前帐号不匹配"})
		return
	}
	message, err1 := Model.GetMessageInfo(fMessageId)
	info, err2 := Model.GetMessageBasicInfo(fMessageId)
	if err1 != nil || err2 != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "服务器异常"})
		return
	}
	if info.MessageType == 3 {
		orderId, _ := strconv.ParseInt(message.Context, 10, 64)
		orderInfo, _ := Model.GetOrderInfo(orderId)
		data := gin.H{
			"data":   orderInfo,
			"status": Model.CheckOrderStatus(orderId, "议价"),
		}
		c.JSON(http.StatusOK, data)
	} else {
		c.JSON(http.StatusOK, message)
	}
}

func DeleteMessage(c *gin.Context) {
	companyId, _ := c.Get("companyId")
	messageId := c.Query("messageId")
	mid, _ := strconv.ParseInt(messageId, 10, 64)
	if !Model.CheckMessageAuth(mid, companyId.(int64)) {
		c.JSON(http.StatusBadRequest, gin.H{"message": "信息的接受帐号与当前帐号不匹配"})
		return
	} else if !Model.DeleteMessage(mid) {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "服务器异常"})
		return
	}
	c.JSON(http.StatusOK, nil)
}
