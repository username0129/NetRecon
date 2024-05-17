package controller

import (
	"backend/internal/model/common"
	"github.com/gin-gonic/gin"
	"net/http"
)

type BaseController struct{}

// GetHealth 获取当前服务状态
func (bc *BaseController) GetHealth(c *gin.Context) {
	common.ResponseOk(c, http.StatusOK, "服务运行正常", nil)
}
