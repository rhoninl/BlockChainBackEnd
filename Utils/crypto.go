package Utils

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
)

const key = "sunyuqingcnmlgcb"

func ParsePassword(password string) (string, error) {
	defer func() {
		recover()
	}()
	myKey := []byte(key)
	myPassword, _ := base64.StdEncoding.DecodeString(password)
	block, _ := aes.NewCipher(myKey)
	blockSize := block.BlockSize()
	blockMode := cipher.NewCBCDecrypter(block, myKey[:blockSize])
	origData := make([]byte, len(myPassword))
	blockMode.CryptBlocks(origData, myPassword)
	origData = PKCS7UnPadding(origData)
	return string(origData), nil
}

func PKCS7UnPadding(origData []byte) []byte {
	length := len(origData)
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}
