package controller

import (
	"backend/internal/global"
	"backend/internal/model"
	"backend/internal/model/common"
	"backend/internal/model/request"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
)

type TaskController struct {
	JWTRequired bool
}

func (tc *TaskController) PostCancelTask(c *gin.Context) {
	var cancelTaskRequest request.CancelTaskRequest

	if err := c.ShouldBindJSON(&cancelTaskRequest); err != nil {
		global.Logger.Error("PostCancelTask 参数解析错误: ", zap.Error(err))
		common.Response(c, http.StatusBadRequest, "参数解析错误", nil)
		return
	}

	if cancel, exists := global.TaskManager[cancelTaskRequest.UUID]; exists {
		var task model.Task
		if err := global.DB.Model(&model.Task{}).Where(&model.Task{UUID: cancelTaskRequest.UUID}).Updates(&task).Error; err != nil {
			global.Logger.Error("未找到对应任务")
			common.Response(c, http.StatusNotFound, "未找到对应任务", nil)
			return
		} else {
			cancel()
			// 取消任务
			if err := task.UpdateStatus("已取消"); err != nil { // 更新任务状态
				global.Logger.Error("更新任务状态失败", zap.String("UUID", cancelTaskRequest.UUID.String()), zap.Error(err))
				common.Response(c, http.StatusInternalServerError, "更新任务状态失败", nil)
				return
			} else {
				common.Response(c, http.StatusOK, "更新任务状态成功", nil)
				return
			}
		}
	} else {
		global.Logger.Error("未找到对应任务")
		common.Response(c, http.StatusNotFound, "未找到对应任务", nil)
		return
	}
}
