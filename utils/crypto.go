package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"template/config"
)

// createKeyFromSecret 使用 SHA-256 从 AppSecret 派生一个 32 字节的密钥
func createKeyFromSecret() []byte {
	// 使用 SHA-256 将任意字符串哈希为 32 字节的 AES-256 密钥
	hash := sha256.Sum256([]byte(config.Config.AppSecret))
	// hash 是一个 [32]byte 数组, 我们返回它的切片
	return hash[:]
}

// Encrypt 使用包级别的 AppSecret 来加密明文字符串
// 返回 Base64 编码的密文字符串
func Encrypt(plaintext string) (string, error) {
	// 1. 派生密钥
	key := createKeyFromSecret()
	plaintextBytes := []byte(plaintext)

	// 2. 创建 AES 密码块
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	// 3. 创建 GCM 
	// AES-GCM 是推荐的带认证的加密模式
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	// 4. 创建一个随机的 Nonce (Number used once)
	// GCM 推荐 12 字节的 Nonce
	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	// 5. 加密数据
	// gcm.Seal 会将 nonce 预置(prepend)到密文的开头
	// 格式为: [nonce][encrypted_data]
	ciphertextBytes := gcm.Seal(nonce, nonce, plaintextBytes, nil)

	// 6. 将二进制密文编码为 Base64 字符串
	// URLEncoding 确保字符串可以安全地用于 URL 和文件名
	return base64.URLEncoding.EncodeToString(ciphertextBytes), nil
}

// Decrypt 使用包级别的 AppSecret 来解密 Base64 编码的密文字符串
func Decrypt(ciphertext string) (string, error) {
	// 1. 派生密钥 (必须与加密时完全相同)
	key := createKeyFromSecret()

	// 2. 将 Base64 字符串解码回二进制字节
	ciphertextBytes, err := base64.URLEncoding.DecodeString(ciphertext)
	if err != nil {
		return "", fmt.Errorf("base64解码失败: %w", err)
	}

	// 3. 创建 AES 密码块
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	// 4. 创建 GCM
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	// 5. 分离 Nonce 和 实际的密文
	nonceSize := gcm.NonceSize()
	if len(ciphertextBytes) < nonceSize {
		return "", errors.New("密文太短，无法分离nonce")
	}

	nonce, actualCiphertext := ciphertextBytes[:nonceSize], ciphertextBytes[nonceSize:]

	// 6. 解密
	// gcm.Open 会自动检查 Nonce、密文和认证标签
	// 如果密钥错误，或者数据在传输过程中被篡改，这里会返回一个 error
	plaintextBytes, err := gcm.Open(nil, nonce, actualCiphertext, nil)
	if err != nil {
		// 这是最常见的错误：密钥不匹配或数据损坏
		return "", fmt.Errorf("解密失败 (密钥错误或数据被篡改): %w", err)
	}

	// 7. 将解密后的字节转换回字符串
	return string(plaintextBytes), nil
}