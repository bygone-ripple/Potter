package utils

import (
	"errors"
	"strings"

	"github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
)

// IsDuplicateKeyError 检查是否为重复键错误
// 兼容不同数据库驱动的重复键错误
func IsDuplicateKeyError(err error) bool {
	if err == nil {
		return false
	}

	// 先检查 GORM 的通用错误
	if errors.Is(err, gorm.ErrDuplicatedKey) {
		return true
	}

	// 检查 MySQL 错误码 1062 (Duplicate entry)
	var mysqlErr *mysql.MySQLError
	if errors.As(err, &mysqlErr) && mysqlErr.Number == 1062 {
		return true
	}

	// 备用方案：检查错误消息中是否包含重复键相关字符串
	errMsg := strings.ToLower(err.Error())
	return strings.Contains(errMsg, "duplicate entry") ||
		strings.Contains(errMsg, "duplicate key")
}
