package controller

import (
	"backend/internal/global"
	"backend/internal/model"
	"backend/internal/model/response"
	"backend/internal/service"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
)

type InitController struct{}

var initService = new(service.InitService)

// PostInit
//
//	@Description: 初始化数据库
//	@receiver ic
//	@param c
//	@Router: /init/init
func (ic *InitController) PostInit(c *gin.Context) {
	if global.DB != nil {
		global.Logger.Error("已存在数据库配置")
		response.Response(c, http.StatusInternalServerError, "已存在数据库配置", nil)
		return
	}

	var req model.InitRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		global.Logger.Error("参数解析错误", zap.Error(err))
		response.Response(c, http.StatusInternalServerError, "参数解析错误", nil)
		return
	}

	if err := initService.Init(req); err != nil {
		global.Logger.Error("数据库初始化错误", zap.Error(err))
		response.Response(c, http.StatusInternalServerError, "数据库初始化错误", nil)
		return
	}
	response.Response(c, http.StatusOK, "数据库初始化成功", nil)
}
