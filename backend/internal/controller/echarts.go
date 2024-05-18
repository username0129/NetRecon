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

// PostFetchTaskCount 获取任务个数
func (ec *EchartsController) PostFetchTaskCount(c *gin.Context) {
	result, err := service.EchartsServiceApp.FetchTasksCount(global.DB, util.GetUUID(c), util.GetAuthorityId(c))
	if err != nil {
		common.ResponseOk(c, http.StatusInternalServerError, "查询数据失败", nil)
		return
	}
	common.ResponseOk(c, http.StatusOK, "查询成功", result)
	return
}

// PostFetchPortCount 获取资产端口数量
func (ec *EchartsController) PostFetchPortCount(c *gin.Context) {
	result, err := service.EchartsServiceApp.FetchPortCount(global.DB, util.GetUUID(c), util.GetAuthorityId(c))
	if err != nil {
		common.ResponseOk(c, http.StatusInternalServerError, "查询数据失败", nil)
		return
	}
	if len(result) == 0 {
		common.ResponseOk(c, http.StatusNotFound, "查无数据", nil)
		return
	}
	common.ResponseOk(c, http.StatusOK, "查询成功", result)
	return
}

// PostFetchDomainCount 获取资产域名数量
func (ec *EchartsController) PostFetchDomainCount(c *gin.Context) {
	result, err := service.EchartsServiceApp.FetchDomainCount(global.DB, util.GetUUID(c), util.GetAuthorityId(c))
	if err != nil {
		common.ResponseOk(c, http.StatusInternalServerError, "查询数据失败", nil)
		return
	}
	if len(result) == 0 {
		common.ResponseOk(c, http.StatusNotFound, "查无数据", nil)
		return
	}
	common.ResponseOk(c, http.StatusOK, "查询成功", result)
	return
}
