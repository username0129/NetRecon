package controller

import (
	"backend/internal/global"
	"backend/internal/model/common"
	"backend/internal/model/request"
	"backend/internal/model/response"
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

	err := service.PortServiceApp.ExecutePortScan(c, portScanRequest, util.GetUUID(c))
	if err != nil {
		common.Response(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}
	common.Response(c, http.StatusOK, "任务提交成功", nil)
	return
}

func (pc *PortScanController) PostFetchResult(c *gin.Context) {
	var req request.FetchResultRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		global.Logger.Error("PostSearchPortScanResult 参数解析错误: ", zap.Error(err))
		common.Response(c, http.StatusBadRequest, "参数解析错误", nil)
		return
	}

	result, total, err := service.PortServiceApp.FetchResult(global.DB, req.PortScanResult, req.PageInfo, req.OrderKey, req.Desc)
	if err != nil {
		global.Logger.Error("查询数据失败: ", zap.Error(err))
		common.Response(c, http.StatusInternalServerError, "查询数据失败", nil)
		return
	}

	if total == 0 {
		common.Response(c, http.StatusNotFound, "未查询到有效数据", nil)
		return
	} else {
		common.Response(c, http.StatusOK, "查询数据成功", response.PageResult{
			Data:     result,
			Total:    total,
			Page:     req.Page,
			PageSize: req.PageSize,
		})
		return
	}
}
