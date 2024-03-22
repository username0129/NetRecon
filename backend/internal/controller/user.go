package controller

import (
	"backend/internal/model/response"
	"github.com/gin-gonic/gin"
	"net/http"
)

type UserController struct {
	JWTRequired bool
}

func (uc *UserController) PostUserInfo(c *gin.Context) {
	response.Response(c, http.StatusOK, "", gin.H{
		"message": "获取用户信息",
	})
}
