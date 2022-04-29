package Utils

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
)

const key = "sunyuqingcnmlgcb"

func AesDecryptCBC(data string) (string, error) {
	defer func() {
		recover()
	}()
	myKey := []byte(key)
	myPassword, _ := base64.StdEncoding.DecodeString(data)
	block, _ := aes.NewCipher(myKey)
	blockSize := block.BlockSize()
	blockMode := cipher.NewCBCDecrypter(block, myKey[:blockSize])
	origData := make([]byte, len(myPassword))
	blockMode.CryptBlocks(origData, myPassword)
	origData = PKCS7UnPadding(origData)
	return string(origData), nil
}

func AesEncryptCBC(data string) string {
	// 分组秘钥
	// NewCipher该函数限制了输入k的长度必须为16, 24或者32
	origData := []byte(data)
	myKey := []byte(key)
	block, _ := aes.NewCipher(myKey)
	blockSize := block.BlockSize()                                // 获取秘钥块的长度
	origData = pkcs7Padding(origData, blockSize)                  // 补全码
	blockMode := cipher.NewCBCEncrypter(block, myKey[:blockSize]) // 加密模式
	encrypted := make([]byte, len(origData))                      // 创建数组
	blockMode.CryptBlocks(encrypted, origData)                    // 加密
	return base64.StdEncoding.EncodeToString(encrypted)
}

func pkcs7Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padText := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padText...)
}

func PKCS7UnPadding(origData []byte) []byte {
	length := len(origData)
	unPadding := int(origData[length-1])
	return origData[:(length - unPadding)]
}
