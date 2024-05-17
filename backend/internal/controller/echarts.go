package controller

import (
	"backend/internal/global"
	"backend/internal/model/common"
	"backend/internal/service"
	"backend/internal/util"
	"github.com/gin-gonic/gin"
	"net/http"
)

type EchartsController struct {
	JWTRequired bool
}

func (ec *EchartsController) PostFetchTaskCount(c *gin.Context) {
	result, err := service.EchartsServiceApp.FetchTasksCount(global.DB, util.GetUUID(c))
	if err != nil {
		common.ResponseOk(c, http.StatusInternalServerError, "查询数据失败", nil)
		return
	}
	common.ResponseOk(c, http.StatusOK, "查询成功", result)
	return
}

func (ec *EchartsController) PostFetchResultCount(c *gin.Context) {
	result, err := service.EchartsServiceApp.FetchTasksCount(global.DB, util.GetUUID(c))
	if err != nil {
		common.ResponseOk(c, http.StatusInternalServerError, "查询数据失败", nil)
		return
	}
	common.ResponseOk(c, http.StatusOK, "查询成功", result)
	return
}
