package Utils

import (
	"math/rand"
	"net/smtp"
	"strconv"
	"strings"
	"time"
)

func GeneratePassWord() string {
	rand.Seed(time.Now().Unix())
	a := rand.Intn(99999)
	b := rand.Intn(9999)
	result := string((a<<5^b)<<2%26+'A') + string((a<<8|b)>>1%26+'A')
	result += FormatString(b, 4) + FormatString(a, 5)
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

func SendMessage(message, email string) error {
	user := `dmutreehole@163.com`
	password := `DLCHYHPHXZVTIIGJ`
	host := `smtp.163.com:25`
	to := email
	subject := `ChainBlock`
	body := message
	err := sendToMail(user, password, host, to, subject, body, "html")
	if err != nil {
		return err
	}
	return nil
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
	if strings.ToLower(code) == strings.ToLower(email.Code) {
		RDB().Del(email.ToEmail + "#emailCode")
		return true, nil
	}
	return false, nil
}
