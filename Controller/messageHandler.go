package Controller

import (
	"github.com/gin-gonic/gin"
	"main/Model"
	"net/http"
	"strconv"
)

func GetAllMessage(c *gin.Context) {
	companyId, _ := c.Get("companyId")
	message, err := Model.GetAllMessage(companyId.(int64))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "服务器异常"})
		return
	}
	c.JSON(http.StatusOK, message)
}

func GetMessageInfo(c *gin.Context) {
	messageId := c.Param("id")
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
	if !Model.CheckMessageAuth(companyId.(int64), fMessageId) {
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
