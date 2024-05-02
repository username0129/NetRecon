package controller

import (
	"backend/internal/global"
	"backend/internal/model/common"
	"backend/internal/service"
	"backend/internal/util"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
)

type UserController struct {
	JWTRequired bool
}

func (uc *UserController) GetUserInfo(c *gin.Context) {
	uuid := util.GetUUID(c)
	if user, err := service.UserServiceApp.GetUserInfo(uuid); err != nil {
		global.Logger.Error("获取用户信息失败: ", zap.Error(err))
		common.Response(c, http.StatusInternalServerError, "获取用户信息失败", nil)
		return
	} else {
		common.Response(c, http.StatusOK, "获取用户信息成功", user)
	}
}

func (uc *UserController) PostUserInfo(c *gin.Context) {
	common.Response(c, http.StatusOK, "", gin.H{
		"message": "获取用户信息",
	})
}
