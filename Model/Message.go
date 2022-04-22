package Model

import (
	"errors"
	"log"
	"main/Utils"
	"sync"
)

func SendMessageTo(messageType int, message string, toId int64, fromId int64) bool {
	go UseClient().SendMessageToId("有一条新消息", toId)
	template := `Insert Into MessageQueue Set MessageType = ? , FromId = ?,ToId = ?,SendTime=now()`
	result, err := Utils.DB().Exec(template, messageType, fromId, toId)
	if err != nil {
		log.Panicln("[SendMessageTo]服务器异常")
		return false
	}
	messageId, _ := result.LastInsertId()
	template = `Insert Into MessageInfo Set MessageId = ?,MessageContent = ?`
	result, err = Utils.DB().Exec(template, messageId, message)
	if err != nil {
		log.Panicln("[SendMessageTo]服务器异常")
		return false
	}
	return true
}

func GetAllMessage(CompanyId int64) ([]Utils.MessageList, error) {
	template := `Select MessageId, MessageType, FromId, isRead,SendTime From MessageQueue Where ToId = ? And isDelete = 0`
	rows, err := Utils.DB().Query(template, CompanyId)
	if err != nil {
		log.Println("[GetAllMessage]服务器异常")
		return nil, err
	}
	defer rows.Close()
	var messageList []Utils.MessageList
	var message Utils.MessageList
	var companyId int64
	wg := sync.WaitGroup{}
	for rows.Next() {
		wg.Add(1)
		rows.Scan(&message.MessageId, &message.MessageType, &companyId, &message.IsRead, &message.SendTime)
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
	if !rows.Next() {
		return false
	}
	var cid int64
	rows.Scan(&cid)
	return cid == CompanyId
}
