package controller

import (
	"backend/internal/global"
	"backend/internal/model/common"
	"backend/internal/model/request"
	"backend/internal/service"
	"backend/internal/util"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
)

type PortScanController struct {
	JWTRequired bool
}

// PostPortScan PostPort 执行端口扫描
func (pc *PortScanController) PostPortScan(c *gin.Context) {
	var portScanRequest request.PortScanRequest

	if err := c.ShouldBindJSON(&portScanRequest); err != nil {
		global.Logger.Error("PostPortScan 参数解析错误: ", zap.Error(err))
		common.Response(c, http.StatusBadRequest, "参数解析错误", nil)
		return
	}

	err := service.PortServiceApp.ExecutePortScan(portScanRequest, util.GetUUID(c))
	if err != nil {
		common.Response(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}
	common.Response(c, http.StatusOK, "任务提交成功", nil)
	return
}
