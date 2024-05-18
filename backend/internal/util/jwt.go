package util

import (
	"backend/internal/global"
	"backend/internal/model"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"net"
	"time"
)

var (
	signingKey = []byte(global.Config.Jwt.SigningKey) // 签发 Token 密钥
)

func GenerateJWT(c model.CustomClaims) (string, error) {
	bf, _ := time.ParseDuration(global.Config.Jwt.BufferTime)     // 缓存时间
	ep, _ := time.ParseDuration(global.Config.Jwt.ExpirationTime) // 过期时间

	c.BufferTime = int64(bf / time.Second)
	c.RegisteredClaims = jwt.RegisteredClaims{
		NotBefore: jwt.NewNumericDate(time.Now().Add(-1000)), // 令牌生效时间
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(ep)),    // 令牌过期时间
		Issuer:    global.Config.Jwt.Issuer,                  // 令牌发行者
	}

	// 使用HS256签名算法
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, c)
	return t.SignedString(signingKey)
}

func ParseToken(tokenString string) (*model.CustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &model.CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return signingKey, nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(*model.CustomClaims); ok && token.Valid {
		return claims, nil
	} else {
		return nil, fmt.Errorf("令牌无效")
	}
}

func SetToken(c *gin.Context, token string, maxAge int) {
	// 增加cookie x-token 向来源的web添加
	host, _, err := net.SplitHostPort(c.Request.Host)
	if err != nil {
		host = c.Request.Host
	}

	if net.ParseIP(host) != nil {
		c.SetCookie("x-token", token, maxAge, "/", "", false, false)
	} else {
		c.SetCookie("x-token", token, maxAge, "/", host, false, false)
	}
}
