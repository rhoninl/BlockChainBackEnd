package Utils

import (
	"github.com/emersion/go-sasl"
	"github.com/emersion/go-smtp"
	"log"
	"math/rand"
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
	auth := sasl.NewPlainClient("", user, password)

	// Connect to the server, authenticate, set the sender and recipient,
	// and send the email all in one step.
	msg := strings.NewReader("To: " + to + "\r\n" +
		"Subject:" + subject + "\r\n" +
		"\r\n" + message +
		".\r\n")
	err := smtp.SendMail(host, auth, user, []string{to}, msg)
	if err != nil {
		log.Fatal(err)
		return err
	}
	return nil
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
