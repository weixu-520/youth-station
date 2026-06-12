package utils

import (
	"github.com/zeromicro/go-zero/core/logx"
	"golang.org/x/crypto/bcrypt"
)

// HashPassword 对密码进行 bcrypt 加密
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		logx.Errorf("failed to hash password: %v", err)
		return "", err
	}
	return string(bytes), nil
}

// CheckPasswordHash 校验密码是否匹配
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {
		logx.Errorf("failed to check password hash: %v", err)
		return false
	}
	return true
}
