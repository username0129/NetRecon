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

type TaskController struct {
	JWTRequired bool
}

// GetAllTasks 获取所有任务列表
func (tc *TaskController) GetAllTasks(c *gin.Context) {
	tasks, err := service.TaskServiceApp.FetchAllTasks()
	if err != nil {
		global.Logger.Error("获取任务列表失败: ", zap.Error(err))
		common.Response(c, http.StatusInternalServerError, "获取任务列表失败", nil)
		return
	}
	if len(tasks) == 0 {
		common.Response(c, http.StatusOK, "任务列表为空", nil)
		return
	} else {
		common.Response(c, http.StatusOK, "获取任务列表成功", tasks)
		return
	}
}

// PostCancelTask 取消指定任务
func (tc *TaskController) PostCancelTask(c *gin.Context) {
	var cancelTaskRequest request.CancelTaskRequest

	if err := c.ShouldBindJSON(&cancelTaskRequest); err != nil {
		global.Logger.Error("PostCancelTask 参数解析错误: ", zap.Error(err))
		common.Response(c, http.StatusBadRequest, "参数解析错误", nil)
		return
	}

	if err := service.TaskServiceApp.CancelTask(cancelTaskRequest.UUID, util.GetUUID(c), util.GetAuthorityId(c)); err != nil {
		global.Logger.Error("CancelTask 运行失败: ", zap.Error(err))
		common.Response(c, http.StatusBadRequest, "CancelTask 运行失败", nil)
		return
	} else {
		common.Response(c, http.StatusBadRequest, "目标任务取消成功", nil)
		return
	}
}
