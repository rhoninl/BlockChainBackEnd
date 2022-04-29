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
		c.JSON(http.StatusBadRequest, gin.H{"message": "滚蛋，消息不是你的"})
		return
	}
	message, err := Model.GetMessageInfo(fMessageId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "服务器异常"})
		return
	}
	c.JSON(http.StatusOK, message)
}

func DeleteMessage(c *gin.Context) {
	companyId, _ := c.Get("companyId")
	messageId := c.Query("messageId")
	mid, _ := strconv.ParseInt(messageId, 10, 64)
	if !Model.CheckMessageAuth(mid, companyId.(int64)) {
		c.JSON(http.StatusBadRequest, gin.H{"message": "消息不是你的删个寄吧"})
		return
	} else if !Model.DeleteMessage(mid) {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "服务器异常"})
		return
	}
	c.JSON(http.StatusOK, nil)
}
