package middleware

import (
	"backend/internal/model/response"
	"backend/internal/util"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"strings"
)

func JWTAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		const BearerSchema = "Bearer "
		authHeader := c.GetHeader("Authorization")
		// 检查是否有请求头
		if authHeader == "" {
			response.Response(c, http.StatusUnauthorized, "用户未登录！", gin.H{"reload": true})
			c.Abort()
			return
		}

		// 检查 Token 是否有 Bearer 前缀
		if !strings.HasPrefix(authHeader, BearerSchema) {
			response.Response(c, http.StatusUnauthorized, "无效的 Token 格式！", gin.H{"reload": true})
			c.Abort()
			return
		}

		// 提取实际的 Token
		tokenString := authHeader[len(BearerSchema):]

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
