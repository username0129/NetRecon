package controller

import (
	"backend/internal/model/response"
	"github.com/gin-gonic/gin"
	"net/http"
)

type BaseController struct{}

// GetHealth
//
//	@Description: 获取当前服务状态
//	@receiver bc
//	@param c
//	@Router: /base/health
func (bc *BaseController) GetHealth(c *gin.Context) {
	response.Response(c, http.StatusOK, "服务运行正常", nil)
}
