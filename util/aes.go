package util

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
)

func PKCS7Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padText := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padText...)
}

func PKCS7UnPadding(origData []byte) []byte {
	length := len(origData)
	unPadding := int(origData[length-1])
	return origData[:(length - unPadding)]
}

//AesEncrypt AES加密,CBC
func AesEncrypt(origData, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	blockSize := block.BlockSize()
	origData = PKCS7Padding(origData, blockSize)
	blockMode := cipher.NewCBCEncrypter(block, key[:blockSize])
	decrypted := make([]byte, len(origData))
	blockMode.CryptBlocks(decrypted, origData)
	return decrypted, nil
}

//AesDecrypt AES解密
func AesDecrypt(encrypted, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	blockSize := block.BlockSize()
	blockMode := cipher.NewCBCDecrypter(block, key[:blockSize])
	origData := make([]byte, len(encrypted))
	blockMode.CryptBlocks(origData, encrypted)
	origData = PKCS7UnPadding(origData)
	return origData, nil
}
