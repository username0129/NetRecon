package middleware

import (
	"backend/internal/model/response"
	"backend/internal/util"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
)

func JWTAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString, err := util.GetToken(c)
		if err != nil {
			response.Response(c, http.StatusUnauthorized, fmt.Sprintf("Token 验证失败: %v", err.Error()), gin.H{"reload": true})
			c.Abort()
			return
		}

		// 解析 Token
		claims, err := util.ParseToken(tokenString)
		if err != nil {
			errorMsg := "令牌无效"
			if errors.Is(err, jwt.ErrTokenExpired) {
				errorMsg = "令牌已过期！"
			}
			response.Response(c, http.StatusUnauthorized, errorMsg, gin.H{"reload": true})
			c.Abort()
			return
		}

		// Token 验证通过，将 claims 保存到请求上下文中
		c.Set("claims", claims)
		c.Next()
	}
}
