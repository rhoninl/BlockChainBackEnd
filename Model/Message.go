package Model

import (
	"errors"
	"github.com/gin-gonic/gin"
	"log"
	"main/Utils"
)

func SendMessageTo(messageType int, message string, toId int64, fromId int64) bool {
	go UseClient().SendMessageToId(gin.H{"message": "有一条新消息", "messageType": 1}, toId)
	template := `Insert Into MessageQueue Set MessageType = ? , FromId = ?,ToId = ?,SendTime=now()`
	result, err := Utils.DB().Exec(template, messageType, fromId, toId)
	if err != nil {
		log.Panicln("[SendMessageTo]服务器异常")
		return false
	}
	messageId, err := result.LastInsertId() ////
	template = `Insert Into MessageInfo Set MessageId = ?,MessageContent = ?`
	result, err = Utils.DB().Exec(template, messageId, message)
	if err != nil {
		log.Panicln("[SendMessageTo]服务器异常")
		return false
	}
	return true
}

func GetMessage(CompanyId, MessageId int64) ([]Utils.MessageList, error) {
	template := `Select MessageId, MessageType, FromId, isRead,SendTime,isReply From MessageQueue Where ToId = ? And isDelete = 0 And MessageId > ? Order By MessageId Desc`
	rows, err := Utils.DB().Query(template, CompanyId, MessageId)
	if err != nil {
		log.Println("[GetMessage]服务器异常")
		return nil, err
	}
	defer rows.Close()
	var messageList []Utils.MessageList
	var message Utils.MessageList
	var companyId int64
	for rows.Next() {
		rows.Scan(&message.MessageId, &message.MessageType, &companyId, &message.IsRead, &message.SendTime, &message.IsReply)
		message.CompanyName, _ = GetCompanyBasicInfo(companyId)
		messageList = append(messageList, message)
	}
	return messageList, nil
}

func GetMessageInfo(MessageId int64) (Utils.Message, error) {
	var message Utils.Message
	template := `Select MessageContent From MessageInfo Where MessageId = ? Limit 1`
	rows, err := Utils.DB().Query(template, MessageId)
	if err != nil {
		log.Println("[GetMessageInfo]", err)
		return message, err
	}
	defer rows.Close()
	if !rows.Next() {
		return message, errors.New("404")
	}
	rows.Scan(&message.Context)
	go func(messageId int64) {
		template1 := `Update MessageQueue Set isRead = 1 Where MessageId = ?`
		Utils.DB().Exec(template1, MessageId)
	}(MessageId)
	return message, nil
}

func CheckMessageAuth(MessageId, CompanyId int64) bool {
	template := `Select ToId From MessageQueue Where MessageId = ?`
	rows, err := Utils.DB().Query(template, MessageId)
	if err != nil {
		return false
	}
	defer rows.Close()
	if !rows.Next() {
		return false
	}
	var cid int64
	rows.Scan(&cid)
	return cid == CompanyId
}

func DeleteMessage(MessageId int64) bool {
	template := `Update MessageQueue Set isDelete = 1 Where MessageId = ?`
	rows, err := Utils.DB().Exec(template, MessageId)
	if err != nil {
		log.Println("[DeleteMessage]Make a mistake")
		return false
	}
	num, _ := rows.RowsAffected()
	return num == 1
}

func GetUnReadNum(CompanyId int64) (int64, error) {
	template := `Select Count(*) From MessageQueue Where ToId = ? And isRead = 0 And isDelete = 0 Group By ToId Limit 1`
	rows, err := Utils.DB().Query(template, CompanyId)
	if err != nil {
		log.Println("[GetUnReadNum]Make a mistake", err)
		return 0, err
	}
	defer rows.Close()
	rows.Next()
	var num int64
	rows.Scan(&num)
	return num, nil
}

func GetMessageBasicInfo(MessageId int64) (Utils.MessageInfo, error) {
	var messageInfo Utils.MessageInfo
	template := `Select MessageType, FromId, ToId ,isReply From MessageQueue Where MessageId = ?`
	rows, err := Utils.DB().Query(template, MessageId)
	if err != nil {
		return messageInfo, err
	}
	defer rows.Close()
	if rows.Next() {
		rows.Scan(&messageInfo.MessageType, &messageInfo.FromId, &messageInfo.ToId, &messageInfo.IsReply)
	}
	return messageInfo, nil
}

func SetReply(MessageId int64) bool {
	template := `Update MessageQueue Set isReply = 1 Where MessageId = ? limit 1`
	result, err := Utils.DB().Exec(template, MessageId)
	if err != nil {
		log.Println(`[SetReply] Make a mistake`)
		return false
	}
	num, _ := result.RowsAffected()
	return num == 1
}
