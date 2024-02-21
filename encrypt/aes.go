// Package encrypt 高级加密标准
// Adevanced Encryption Standard, AES
package encrypt

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"errors"
)

// 16,24,32位字符串的话，分别对应AES-128，AES-192，AES-256 加密方法
var passwordKey = []byte("DIS**#KKKDJJSKDI")

// PKCS7Padding PKCS7 填充模式
func PKCS7Padding(data []byte, blockSize int) []byte {
	padding := blockSize - len(data)%blockSize
	// 把切片 []byte{byte(padding)} 复制padding个
	// 合并成新的字节切片返回
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)

	return append(data, padtext...)
}

// PKCS7UnPadding 填充的逆向操作
func PKCS7UnPadding(data []byte) ([]byte, error) {
	// 获取数据长度
	length := len(data)

	if length == 0 {
		return nil, errors.New("加密字符串错误！")
	}
	// 获取填充字符串长度
	unpadding := int(data[length-1])
	// 截取切片，删除填充字节，并且返回明文
	return data[:(length - unpadding)], nil
}

// AesEncrypt 加密
func AesEncrypt(data, key []byte) ([]byte, error) {
	// 创建加密算法实例
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	// 获取块的大小
	blockSize := block.BlockSize()
	// 对数据进行填充，让数据长度满足需求
	data = PKCS7Padding(data, blockSize)
	// 采用AES加密方法中CBC加密模式
	blocMode := cipher.NewCBCEncrypter(block, key[:blockSize])
	crypted := make([]byte, len(data))
	// 执行加密
	blocMode.CryptBlocks(crypted, data)
	return crypted, nil
}

// AesDecrypt 解密
func AesDecrypt(cypted []byte, key []byte) ([]byte, error) {
	// 创建加密算法实例
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	// 获取块大小
	blockSize := block.BlockSize()
	// 创建加密客户端实例
	blockMode := cipher.NewCBCDecrypter(block, key[:blockSize])
	data := make([]byte, len(cypted))
	// 这个函数也可以用来解密
	blockMode.CryptBlocks(data, cypted)
	// 去除填充字符串
	data, err = PKCS7UnPadding(data)
	if err != nil {
		return nil, err
	}
	return data, err
}

// EnPasswordCode 加密 base64
func EnPasswordCode(password []byte) (string, error) {
	result, err := AesEncrypt(password, passwordKey)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(result), err
}

// DePasswordCode 解密 base64
func DePasswordCode(password string) ([]byte, error) {
	// 解密base64字符串
	passwordByte, err := base64.StdEncoding.DecodeString(password)
	if err != nil {
		return nil, err
	}
	// 执行AES解密
	return AesDecrypt(passwordByte, passwordKey)
}
