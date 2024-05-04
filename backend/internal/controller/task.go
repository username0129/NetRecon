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

type TaskController struct {
	JWTRequired bool
}

// PostCancelTask 取消指定任务
func (tc *TaskController) PostCancelTask(c *gin.Context) {
	var req request.CancelTaskRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		global.Logger.Error("PostCancelTask 参数解析错误: ", zap.Error(err))
		common.Response(c, http.StatusBadRequest, "参数解析错误", nil)
		return
	}

	if err := service.TaskServiceApp.CancelTask(req.UUID, util.GetUUID(c), util.GetAuthorityId(c)); err != nil {
		global.Logger.Error("CancelTask 运行失败: ", zap.Error(err))
		common.Response(c, http.StatusInternalServerError, "CancelTask 运行失败", nil)
		return
	} else {
		common.Response(c, http.StatusOK, "目标任务取消成功", nil)
		return
	}
}

// PostFetchTasks 获取任务列表
func (tc *TaskController) PostFetchTasks(c *gin.Context) {
	var req request.FetchTasksRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		global.Logger.Error("PostFetchTasks 参数解析错误: ", zap.Error(err))
		common.Response(c, http.StatusBadRequest, "参数解析错误", nil)
		return
	}

	result, total, err := service.TaskServiceApp.FetchTasks(global.DB, req.Task, req.PageInfo, req.OrderKey, req.Desc, util.GetUUID(c), util.GetAuthorityId(c))
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

// PostDeleteTask 删除指定任务及其结果
func (tc *TaskController) PostDeleteTask(c *gin.Context) {
	var req request.DeleteTaskRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		global.Logger.Error("PostCancelTask 参数解析错误: ", zap.Error(err))
		common.Response(c, http.StatusBadRequest, "参数解析错误", nil)
		return
	}

	if err := service.TaskServiceApp.DeleteTask(global.DB, req.UUID, util.GetUUID(c), util.GetAuthorityId(c)); err != nil {
		global.Logger.Error("PostDeleteTask 运行失败: ", zap.Error(err))
		common.Response(c, http.StatusInternalServerError, "删除任务失败", nil)
		return
	} else {
		common.Response(c, http.StatusOK, "目标任务及结果删除成功", nil)
		return
	}
}
