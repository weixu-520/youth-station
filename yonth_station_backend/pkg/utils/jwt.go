package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/zeromicro/go-zero/core/logx"
)

// GenerateToken 生成 JWT token（使用 MapClaims，兼容 go-zero JWT 中间件）
func GenerateToken(userId int64, isAdmin bool, secret string, expireSeconds int64) (string, int64, error) {
	expireAt := time.Now().Unix() + expireSeconds
	claims := jwt.MapClaims{
		"userId":  userId,
		"isAdmin": isAdmin,
		"exp":     expireAt,
		"iat":     time.Now().Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(secret))
	if err != nil {
		logx.Errorf("failed to sign JWT token: %v", err)
		return "", 0, err
	}
	return signedToken, expireAt, nil
}

func ParseToken(tokenString, secret string) (int64, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return 0, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secret), nil
	})
	if err != nil {
		return 0, err
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		userIdVal, ok := claims["userId"]
		if !ok {
			return 0, errors.New("userId claim not found")
		}
		switch v := userIdVal.(type) {
		case float64:
			return int64(v), nil
		case int64:
			return v, nil
		case json.Number:
			return v.Int64()
		default:
			return 0, errors.New("invalid userId type")
		}
	}
	return 0, errors.New("invalid token")
}
