package tools

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"fmt"
)



func pkcs7Unpad(data []byte, blockSize int) ([]byte) {
	if blockSize <= 0 {
		return nil
	}
	if len(data)%blockSize != 0 || len(data) == 0 {
		return nil
	}
	c := data[len(data)-1]
	n := int(c)
	if n == 0 || n > len(data) {
		return nil
	}
	for i := 0; i < n; i++ {
		if data[len(data)-n+i] != c {
			return nil
		}
	}
	return data[:len(data)-n]
}
func AesEncrypt(origDataStr, keyStr string) ([]byte, error) {
	origData:=[]byte(origDataStr)
	key:=[]byte(keyStr)
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	blockSize := block.BlockSize()
	origData = PKCS5Padding(origData, blockSize)
	blockMode := cipher.NewCBCEncrypter(block, key[:blockSize])
	crypted := make([]byte, len(origData))
	blockMode.CryptBlocks(crypted, origData)
	return crypted, nil
}

func AesDecrypt(cryptedStr, keyStr string) ([]byte, error) {
	crypted,_:=base64.StdEncoding.DecodeString(cryptedStr)
	fmt.Println(string(crypted))
	key,_:=base64.StdEncoding.DecodeString(keyStr)
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	blockSize := block.BlockSize()
	fmt.Println("blockSize",blockSize)
	fmt.Println(string(key[:blockSize]))

	blockMode := cipher.NewCBCDecrypter(block, key[:blockSize])
	origData := make([]byte, len(crypted))
	blockMode.CryptBlocks(origData, crypted)
	origData = pkcs7Unpad(origData,blockSize)
	return origData, nil
}