package utils

import (
	"golang.org/x/crypto/bcrypt"
)

// HashPassword 为密码生成一个安全的、加盐的 bcrypt 哈希值。
func HashPassword(password string) (string, error) {
	// bcrypt.GenerateFromPassword 会自动处理：
	// 1. 生成一个随机的盐 (salt)
	// 2. 将盐和密码组合
	// 3. 使用 bcrypt 算法进行哈希
	// 4. 将算法版本、cost、盐 和 哈希值 组合成一个字符串返回

	// bcrypt.DefaultCost (目前是 10) 是一个安全且性能均衡的成本因子。
	// 成本越高，哈希越慢，越安全，但验证也越慢。
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	// 将哈希后的 []byte 转换为 string 存储到数据库
	return string(hashedBytes), nil
}

// CheckPasswordHash 比较用户提交的明文密码 (password) 和数据库中存储的哈希值 (hash)。
func CheckPasswordHash(password string, hash string) bool {
	// bcrypt.CompareHashAndPassword 会：
	// 1. 从 hash 字符串中解析出 盐 (salt) 和 哈希值
	// 2. 使用相同的 盐 和用户提交的 password 进行哈希
	// 3. 比较新生成的哈希和存储的哈希是否一致（使用恒定时间比较）

	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))

	// 如果 err == nil，表示密码匹配
	// 如果 err != nil (例如 bcrypt.ErrMismatchedHashAndPassword)，表示密码不匹配
	return err == nil
}
