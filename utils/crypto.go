package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"errors"
	"fmt"
)

// 定义加密模式（默认CBC，GCM为可选增强模式）
const (
	ModeCBC = "CBC"
	ModeGCM = "GCM"
)

// PKCS7补码（AES块对齐必需）
func pkcs7Padding(data []byte, blockSize int) []byte {
	padding := blockSize - len(data)%blockSize
	padText := make([]byte, padding)
	for i := range padText {
		padText[i] = byte(padding)
	}
	return append(data, padText...)
}

// PKCS7解补码
func pkcs7UnPadding(data []byte) ([]byte, error) {
	length := len(data)
	if length == 0 {
		return nil, errors.New("密文长度不能为0")
	}
	padding := int(data[length-1])
	if padding > length || padding == 0 {
		return nil, errors.New("补码长度不合法")
	}
	return data[:length-padding], nil
}

// AESEncrypt 字符串加密为二进制密文
// 参数：
//
//	plainText: 待加密的字符串
//	key: 加密密钥（长度必须为16/24/32字节，对应AES-128/192/256）
//	iv: 初始化向量（CBC模式需16字节，GCM模式需12字节，nil则自动生成GCM的nonce）
//	mode: 加密模式（ModeCBC/ModeGCM，默认CBC）
//
// 返回：
//
//	cipherData: 加密后的二进制密文
//	usedIV: 实际使用的IV/nonce（GCM模式自动生成时返回）
//	err: 错误信息
func AESEncrypt(plainText string, key, iv []byte, mode string) (cipherData []byte, usedIV []byte, err error) {
	// 1. 校验密钥长度
	keyLen := len(key)
	if keyLen != 16 && keyLen != 24 && keyLen != 32 {
		return nil, nil, fmt.Errorf("密钥长度非法，需16/24/32字节（当前：%d字节）", keyLen)
	}

	// 2. 字符串转字节
	plainData := []byte(plainText)

	// 3. 创建AES加密块
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, nil, fmt.Errorf("创建AES加密块失败：%w", err)
	}

	// 4. 按模式处理加密
	switch mode {
	case ModeGCM:
		// GCM模式（无需补码，自带认证，安全性更高）
		gcm, err := cipher.NewGCM(block)
		if err != nil {
			return nil, nil, fmt.Errorf("初始化GCM模式失败：%w", err)
		}
		// GCM推荐nonce长度为12字节，未传入则自动生成
		if iv == nil {
			usedIV = make([]byte, gcm.NonceSize())
			copy(usedIV, iv[:gcm.NonceSize()])
		} else {
			if len(iv) != gcm.NonceSize() {
				return nil, nil, fmt.Errorf("GCM模式IV长度需为%d字节（当前：%d字节）", gcm.NonceSize(), len(iv))
			}
			usedIV = iv
		}
		// 加密（GCM无需补码）
		cipherData = gcm.Seal(nil, usedIV, plainData, nil)

	case ModeCBC, "":
		// CBC模式（默认，需补码）
		if len(iv) != aes.BlockSize {
			return nil, nil, fmt.Errorf("CBC模式IV长度需为%d字节（当前：%d字节）", aes.BlockSize, len(iv))
		}
		usedIV = iv
		// 补码对齐块大小
		plainData = pkcs7Padding(plainData, aes.BlockSize)
		// 初始化CBC加密器
		modeCBC := cipher.NewCBCEncrypter(block, iv)
		// 加密
		cipherData = make([]byte, len(plainData))
		modeCBC.CryptBlocks(cipherData, plainData)

	default:
		return nil, nil, fmt.Errorf("不支持的加密模式：%s（仅支持%s/%s）", mode, ModeCBC, ModeGCM)
	}

	return cipherData, usedIV, nil
}

// AESDecrypt 二进制密文解密为字符串
// 参数：
//
//	cipherData: 加密后的二进制密文
//	key: 解密密钥（需与加密密钥一致）
//	iv: 加密时使用的IV/nonce（需与加密时一致）
//	mode: 加密模式（需与加密时一致）
//
// 返回：
//
//	plainText: 解密后的字符串
//	err: 错误信息
func AESDecrypt(cipherData []byte, key, iv []byte, mode string) (plainText string, err error) {
	// 1. 校验密钥长度
	keyLen := len(key)
	if keyLen != 16 && keyLen != 24 && keyLen != 32 {
		return "", fmt.Errorf("密钥长度非法，需16/24/32字节（当前：%d字节）", keyLen)
	}

	// 2. 校验密文长度
	if len(cipherData) == 0 {
		return "", errors.New("密文长度不能为0")
	}

	// 3. 创建AES解密块
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", fmt.Errorf("创建AES解密块失败：%w", err)
	}

	// 4. 按模式处理解密
	var plainData []byte
	switch mode {
	case ModeGCM:
		gcm, err := cipher.NewGCM(block)
		if err != nil {
			return "", fmt.Errorf("初始化GCM模式失败：%w", err)
		}
		if len(iv) != gcm.NonceSize() {
			return "", fmt.Errorf("GCM模式IV长度需为%d字节（当前：%d字节）", gcm.NonceSize(), len(iv))
		}
		// 解密（GCM无需解补码）
		plainData, err = gcm.Open(nil, iv, cipherData, nil)
		if err != nil {
			return "", fmt.Errorf("GCM解密失败：%w", err)
		}

	case ModeCBC, "":
		if len(iv) != aes.BlockSize {
			return "", fmt.Errorf("CBC模式IV长度需为%d字节（当前：%d字节）", aes.BlockSize, len(iv))
		}
		// 初始化CBC解密器
		modeCBC := cipher.NewCBCDecrypter(block, iv)
		// 解密
		plainData = make([]byte, len(cipherData))
		modeCBC.CryptBlocks(plainData, cipherData)
		// 解补码
		plainData, err = pkcs7UnPadding(plainData)
		if err != nil {
			return "", fmt.Errorf("CBC解补码失败：%w", err)
		}

	default:
		return "", fmt.Errorf("不支持的解密模式：%s（仅支持%s/%s）", mode, ModeCBC, ModeGCM)
	}

	// 5. 字节转字符串
	return string(plainData), nil
}

// ------------------------------ 便捷封装（可选） ------------------------------
// AESEncryptCBC 简化版CBC加密（仅传key和IV，默认CBC模式）
func AESEncryptCBC(plainText string, key, iv []byte) ([]byte, error) {
	cipherData, _, err := AESEncrypt(plainText, key, iv, ModeCBC)
	return cipherData, err
}

// AESDecryptCBC 简化版CBC解密
func AESDecryptCBC(cipherData []byte, key, iv []byte) (string, error) {
	return AESDecrypt(cipherData, key, iv, ModeCBC)
}
