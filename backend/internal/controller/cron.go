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

type CronController struct {
	JWTRequired bool
}

func (cc *CronController) PostAddTask(c *gin.Context) {
	var req request.CronAddTaskRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		global.Logger.Error("PostAddTask 参数解析错误: ", zap.Error(err))
		common.Response(c, http.StatusBadRequest, "参数解析错误", nil)
		return
	}
	err := service.CronServiceApp.AddTask(global.CronManager, req, util.GetUUID(c), util.GetAuthorityId(c))
	if err != nil {
		common.Response(c, http.StatusInternalServerError, "添加计划任务失败", nil)
		return
	}
	common.Response(c, http.StatusOK, "添加计划任务成功", nil)
	return
}

func (cc *CronController) FetchCronTasks(c *gin.Context) {
	var req request.CronAddTaskRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		global.Logger.Error("PostAddTask 参数解析错误: ", zap.Error(err))
		common.Response(c, http.StatusBadRequest, "参数解析错误", nil)
		return
	}

	return
}
