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

type SubDomainController struct {
	JWTRequired bool
}

func (sc *SubDomainController) PostBruteSubdomains(c *gin.Context) {
	var subDomainRequest request.SubDomainRequest
	if err := c.ShouldBindJSON(&subDomainRequest); err != nil {
		global.Logger.Error("PostBruteSubdomains 参数解析错误: ", zap.Error(err))
		common.Response(c, http.StatusBadRequest, "参数解析错误", nil)
		return
	}

	err := service.SubDomainServiceApp.BruteSubdomains(subDomainRequest, util.GetUUID(c))
	if err != nil {
		common.Response(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}
	common.Response(c, http.StatusOK, "任务提交成功", nil)
	return
}
