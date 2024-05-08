package controller

import (
	"backend/internal/global"
	"backend/internal/model/common"
	"backend/internal/model/request"
	"backend/internal/model/response"
	"backend/internal/service"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
)

type OperationController struct {
	JWTRequired bool
}

func (pc *OperationController) PostFetchResult(c *gin.Context) {
	var req request.FetchOperationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		global.Logger.Error("PostFetchResult 参数解析错误: ", zap.Error(err))
		common.ResponseOk(c, http.StatusBadRequest, "参数解析错误", nil)
		return
	}

	result, total, err := service.OperationServiceApp.FetchResult(global.DB, req.OperationRecord, req.PageInfo, req.OrderKey, req.Desc)
	if err != nil {
		global.Logger.Error("查询数据失败: ", zap.Error(err))
		common.ResponseOk(c, http.StatusInternalServerError, "查询数据失败", nil)
		return
	}

	if total == 0 {
		common.ResponseOk(c, http.StatusNotFound, "未查询到有效数据", nil)
		return
	} else {
		common.ResponseOk(c, http.StatusOK, "查询数据成功", response.PageResult{
			Data:     result,
			Total:    total,
			Page:     req.Page,
			PageSize: req.PageSize,
		})
		return
	}
}

// PostDeleteResult 删除指定记录
func (pc *OperationController) PostDeleteResult(c *gin.Context) {
	var req request.DeleteOperationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		global.Logger.Error("PostDeleteResult 参数解析错误: ", zap.Error(err))
		common.ResponseOk(c, http.StatusBadRequest, "参数解析错误", nil)
		return
	}

	if err := service.OperationServiceApp.DeleteResult(req.UUID); err != nil {
		global.Logger.Error("PostDeleteResult 运行失败: ", zap.Error(err))
		common.ResponseOk(c, http.StatusInternalServerError, "目标结果删除失败", nil)
		return
	} else {
		common.ResponseOk(c, http.StatusOK, "目标结果删除成功", nil)
		return
	}
}
