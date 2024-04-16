package util

import (
	"backend/internal/e"
	"backend/internal/model"
	"github.com/gin-gonic/gin"
	"strings"
)

func GetToken(c *gin.Context) (token string, err error) {
	const BearerSchema = "Bearer "
	authHeader := c.GetHeader("Authorization")
	// 检查是否有请求头
	if authHeader == "" {
		return "", e.ErrUserNotLogin
	}

	// 检查 Token 是否有 Bearer 前缀
	if !strings.HasPrefix(authHeader, BearerSchema) {
		return "", e.ErrTokenMalformed
	}
	// 提取实际的 Token
	token = authHeader[len(BearerSchema):]
	return token, nil
}

func GetClaims(c *gin.Context) *model.CustomClaims {
	token, _ := GetToken(c)
	claims, _ := ParseToken(token)
	return claims
}
