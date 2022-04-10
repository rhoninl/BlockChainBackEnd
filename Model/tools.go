package Model

import (
	"math/rand"
	"strconv"
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
