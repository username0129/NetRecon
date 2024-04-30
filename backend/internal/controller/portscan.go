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

// PostByTaskUUID 执行端口扫描
func (pc *PortScanController) PostByTaskUUID(c *gin.Context) {
	var portScanResultByTaskUUIDRequest request.PortScanResultByTaskUUIDRequest

	if err := c.ShouldBindJSON(&portScanResultByTaskUUIDRequest); err != nil {
		global.Logger.Error("PostPortScan 参数解析错误: ", zap.Error(err))
		common.Response(c, http.StatusBadRequest, "参数解析错误", nil)
		return
	}

	resultList, err := service.PortServiceApp.GetResultsByTaskUUID(global.DB, portScanResultByTaskUUIDRequest.TaskUUID)
	if err != nil {
		common.Response(c, http.StatusInternalServerError, "获取扫描结果失败", nil)
		return
	}
	common.Response(c, http.StatusOK, "获取扫描结果成功", resultList)
	return
}

// PostByTaskUUID 执行端口扫描
func (pc *PortScanController) PostByIP(c *gin.Context) {
	var portScanResultByIPRequest request.PortScanResultByIPRequest

	if err := c.ShouldBindJSON(&portScanResultByIPRequest); err != nil {
		global.Logger.Error("PostPortScan 参数解析错误: ", zap.Error(err))
		common.Response(c, http.StatusBadRequest, "参数解析错误", nil)
		return
	}

	resultList, err := service.PortServiceApp.GetResultsByIP(global.DB, portScanResultByIPRequest.IP)
	if err != nil {
		global.Logger.Error("获取扫描结果失败", zap.Error(err))
		common.Response(c, http.StatusInternalServerError, "获取扫描结果失败", nil)
		return
	}
	common.Response(c, http.StatusOK, "获取扫描结果成功", resultList)
	return
}
