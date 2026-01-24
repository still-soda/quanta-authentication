package utilities

import (
	"crypto"
	"crypto/rand"
	"encoding/hex"
	"strconv"
)

// ParseInt 解析字符串为整数
func ParseInt(s string) (int, error) {
	return strconv.Atoi(s)
}

// ParseIntParam 解析字符串为整数，如果失败则返回默认值
func ParseIntParam(s string, defaultVal int) int {
	if s == "" {
		return defaultVal
	}
	val, err := strconv.Atoi(s)
	if err != nil {
		return defaultVal
	}
	return val
}

// GenerateSalt 生成指定长度的随机盐值
func GenerateSalt(length int) (string, error) {
	b := make([]byte, length)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

// HashPassword 使用盐值对密码进行哈希处理
func HashPassword(password, salt string) string {
	hasher := crypto.SHA256.New()
	hasher.Write([]byte(salt + password))
	return hex.EncodeToString(hasher.Sum(nil))
}

// VerifyPassword 验证输入的密码是否与存储的哈希值匹配
func VerifyPassword(password, salt, hash string) bool {
	hashedInput := HashPassword(password, salt)
	return hashedInput == hash
}
