package controller

import (
	"backend/internal/global"
	"backend/internal/model/common"
	"backend/internal/model/request"
	"backend/internal/model/response"
	"backend/internal/service"
	"errors"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
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

	if req.Username != "" {
		uuid, err := service.UserServiceApp.GetUserUUIDByUsername(req.Username)
		if err != nil {
			global.Logger.Error("查询数据失败: ", zap.Error(err))
			if errors.Is(err, gorm.ErrRecordNotFound) {
				common.ResponseOk(c, http.StatusInternalServerError, "用户名不存在", nil)
			}
			common.ResponseOk(c, http.StatusInternalServerError, "用户名无效", nil)
			return
		}
		req.OperationRecord.UserUUID = uuid
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

// PostDeleteResults 批量删除
func (pc *OperationController) PostDeleteResults(c *gin.Context) {
	var req request.DeleteOperationsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		global.Logger.Error("PostDeleteResult 参数解析错误: ", zap.Error(err))
		common.ResponseOk(c, http.StatusBadRequest, "参数解析错误", nil)
		return
	}

	if err := service.OperationServiceApp.DeleteResults(req.UUIDS); err != nil {
		global.Logger.Error("PostDeleteResult 运行失败: ", zap.Error(err))
		common.ResponseOk(c, http.StatusInternalServerError, "目标结果删除失败", nil)
		return
	} else {
		common.ResponseOk(c, http.StatusOK, "目标结果删除成功", nil)
		return
	}
}
