package crypt

import (
	"bytes"
	"crypto/cipher"
	"crypto/des"
	"crypto/md5"
	"encoding/hex"
	"fmt"
)

// Des-cbc加密
func DesCbcEnc(data []byte, key []byte, iv []byte) ([]byte, error) {
	var err error
	if iv == nil {
		iv = key
	}
	// 创建新的加密块
	block, err := des.NewCipher(key)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	data = Pkcs5Padding(data, block.BlockSize())
	mode := cipher.NewCBCEncrypter(block, iv)
	dst := make([]byte, len(data))
	mode.CryptBlocks(dst, data)
	return dst, nil
}

// Des-cbc 解密
func DesCbcDec(data []byte, key []byte, iv []byte) ([]byte, error) {
	if iv == nil {
		iv = key
	}
	block, err := des.NewCipher(key)
	if err != nil {
		return nil, err
	}
	mode := cipher.NewCBCDecrypter(block, iv)
	plainText := make([]byte, len(data))
	mode.CryptBlocks(plainText, data)
	plainText = PKCS5UnPadding(plainText)
	return plainText, nil
}

func Pkcs5Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padText := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padText...)
}

func PKCS5UnPadding(text []byte) []byte {
	length := len(text)
	unpaddingLen := int(text[length - 1])
	return text[:(length - unpaddingLen)]
}

func SimpleMd5(data string) string {
	h := md5.New()
	h.Write([]byte(data))
	return hex.EncodeToString(h.Sum(nil))
}