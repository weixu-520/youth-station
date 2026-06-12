package utils

import (
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/zeromicro/go-zero/core/logx"
)

// GenerateToken 生成 JWT token（使用 MapClaims，兼容 go-zero JWT 中间件）
func GenerateToken(userId int64, secret string, expireSeconds int64) (string, int64, error) {
	expireAt := time.Now().Unix() + expireSeconds
	claims := jwt.MapClaims{
		"userId": userId,
		"exp":    expireAt,
		"iat":    time.Now().Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(secret))
	if err != nil {
		logx.Errorf("failed to sign JWT token: %v", err)
		return "", 0, err
	}
	return signedToken, expireAt, nil
}
