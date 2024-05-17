package middleware

import (
	"backend/internal/model/common"
	"backend/internal/service"
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
			common.ResponseOk(c, http.StatusUnauthorized, fmt.Sprintf("Token 验证失败: %v", err.Error()), nil)
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
			common.ResponseOk(c, http.StatusUnauthorized, errorMsg, nil)
			c.Abort()
			return
		}
		// 判断 Token 所属用户是否存在
		if _, err = service.UserServiceApp.FetchUserByUUID(claims.UUID); err != nil {
			common.ResponseOk(c, http.StatusUnauthorized, fmt.Sprintf("用户信息被删除或不存在: %v", err.Error()), nil)
			c.Abort()
			return
		}
		c.Next()
	}
}
