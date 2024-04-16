package controller

import (
	"backend/internal/global"
	"backend/internal/model/response"
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
	if routes, err := service.RouterServiceApp.GetRouteTree(util.GetClaims(c).AuthorityId); err != nil {
		global.Logger.Error("获取失败路由失败: ", zap.Error(err))
		response.Response(c, http.StatusInternalServerError, "获取用户路由信息失败", nil)
		return
	} else {
		response.Response(c, http.StatusOK, "获取用户路由信息成功", routes)
		return
	}
}
