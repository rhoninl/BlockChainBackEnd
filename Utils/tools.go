package Utils

import (
	"math/rand"
	"net/smtp"
	"strconv"
	"strings"
	"time"
)

func GenerateId() string {
	rand.Seed(time.Now().Unix())
	a := rand.Intn(9999)
	b := rand.Intn(999)
	result := string((a<<5^b)<<2%26+'A') + string((a<<8|b)>>1%26+'A')
	result += `-` + FormatString(b, 3)
	result += `-` + FormatString(a, 4)
	return result
}

//FormatString 生成制定长度的字符串
func FormatString(i int, length int) string {
	result := strconv.Itoa(i)
	if len(result) > length {
		return result[length:]
	}
	for len(result) < length {
		result = "0" + result
	}
	return result
}

//获取六位随机验证码
func GenVerCode() string {
	result := ""
	directory := `0123456789ABCDEFGHJKLMNPQRSTUVWXYZ`
	for i := 0; i < 6; i++ {
		rand.Seed(time.Now().UnixNano())
		a := rand.Int() % 34
		result += directory[a : a+1]
	}
	return result
}

func SendCode(info AuthCode) error {
	user := `dmutreehole@163.com`
	password := `DLCHYHPHXZVTIIGJ`
	host := `smtp.163.com:25`
	to := info.ToEmail
	subject := `ChainBlock`
	body := `<html><body><a>您的验证码为</a><h3>` + info.Code + `</h3><a><br/>验证码有效期为1小时，请在1小时内完成验证<br/>如果不是您本人操作，请忽略本条邮件</a></body></html>`
	err := sendToMail(user, password, host, to, subject, body, "html")
	if err != nil {
		return err
	}
	return AuthCodeRegister(info)
}

//SendToMail 发送邮件的函数
func sendToMail(user, password, host, to, subject, body, mailType string) error {
	hp := strings.Split(host, ":")
	auth := smtp.PlainAuth("", user, password, hp[0])
	var contentType string
	if mailType == "html" {
		contentType = "Content-Type: text/" + mailType + "; charset=UTF-8"
	} else {
		contentType = "Content-Type: text/plain" + "; charset=UTF-8"
	}
	msg := []byte("To: " + to + "\r\nFrom: " + user + ">\r\nSubject: " + subject + "\r\n" + contentType + "\r\n\r\n" + body)
	sendTo := strings.Split(to, ";")
	err := smtp.SendMail(host, auth, user, sendTo, msg)
	return err
}

func AuthCodeRegister(email AuthCode) error {
	_, err := RDB().Set(email.ToEmail+"#emailCode", email.Code, time.Hour).Result()
	return err
}

func AuthCodeCheck(email AuthCode) (bool, error) {
	code, err := RDB().Get(email.ToEmail + "#emailCode").Result()
	if err != nil {
		return false, err
	}
	return code == email.Code, nil
}
