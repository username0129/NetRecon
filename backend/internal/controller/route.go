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

type RouteController struct {
	JWTRequired bool
}

func (rc *RouteController) GetRoute(c *gin.Context) {
	if routes, err := service.RouterServiceApp.GetRouteTree(util.GetAuthorityId(c)); err != nil {
		global.Logger.Error("获取用户路由信息失败: ", zap.Error(err))
		common.ResponseOk(c, http.StatusInternalServerError, "获取用户路由信息失败", nil)
		return
	} else {
		common.ResponseOk(c, http.StatusOK, "获取用户路由信息成功", routes)
		return
	}
}
