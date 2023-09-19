package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"errors"
)

func AesEncryptByGCM(data, key string) (string, error) {
	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		//panic(fmt.Sprintf("NewCipher error:%s", err))
		return "", err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		//panic(fmt.Sprintf("NewGCM error:%s", err))
		return "", err
	}
	nonceStr := key[:gcm.NonceSize()]
	nonce := []byte(nonceStr)
	seal := gcm.Seal(nonce, nonce, []byte(data), nil)
	return base64.StdEncoding.EncodeToString(seal), nil
}

func AesDecryptByGCM(data, key string) (string, error) {
	// 反解base64
	dataByte, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		//panic(fmt.Sprintf("base64 DecodeString error:%s", err))
		return "", err
	}
	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		//panic(fmt.Sprintf("NewCipher error:%s", err))
		return "", err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		//panic(fmt.Sprintf("NewGCM error:%s", err))
		return "", err
	}
	nonceSize := gcm.NonceSize()
	if len(dataByte) < nonceSize {
		//panic("dataByte to short")
		customErr := errors.New("dataByte to short")
		return "", customErr
	}
	nonce, ciphertext := dataByte[:nonceSize], dataByte[nonceSize:]
	open, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		//panic(fmt.Sprintf("gcm Open error:%s", err))
		return "", err
	}
	return string(open), nil
}
