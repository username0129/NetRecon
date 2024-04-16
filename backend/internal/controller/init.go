package controller

import (
	"backend/internal/global"
	"backend/internal/model"
	"backend/internal/model/response"
	"backend/internal/service"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

type InitController struct{}

// GetInit
//
//	@Description: 检查数据库初始化状态
//	@receiver ic
//	@param c
//	@Router: /init/init
func (ic *InitController) GetInit(c *gin.Context) {
	if global.DB != nil {
		response.Response(c, http.StatusOK, "已存在数据库配置", nil)
		return
	}
	response.Response(c, http.StatusInternalServerError, "数据库尚未初始化", nil)
	return
}

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
		global.Logger.Error(fmt.Sprintf("参数解析错误 %v", err.Error()))
		response.Response(c, http.StatusInternalServerError, "参数解析错误", nil)
		return
	}

	if err := service.InitServiceApp.Init(req); err != nil {
		global.Logger.Error(fmt.Sprintf("数据库初始化错误：%v", err.Error()))
		response.Response(c, http.StatusInternalServerError, "数据库初始化错误，详情请查看后端。", nil)
		return
	}
	response.Response(c, http.StatusOK, "数据库初始化成功", nil)
	return
}
