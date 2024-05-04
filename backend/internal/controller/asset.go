package controller

import (
	"backend/internal/global"
	"backend/internal/model/common"
	"backend/internal/model/request"
	"backend/internal/service"
	"backend/internal/util"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
)

type AssetController struct{}

func (ac *AssetController) PostAddAsset(c *gin.Context) {
	var req request.AddAssetRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		global.Logger.Error("PostAddAsset 参数解析错误: ", zap.Error(err))
		common.Response(c, http.StatusBadRequest, "参数解析错误", nil)
		return
	}
	if err := service.AssetServiceApp.AddAsset(global.DB, req, util.GetUUID(c)); err != nil {
		global.Logger.Error("添加资产错误: ", zap.Error(err))
		common.Response(c, http.StatusBadRequest, fmt.Sprintf("添加资产错误：%v", err.Error()), nil)
		return
	}
	common.Response(c, http.StatusBadRequest, "添加资产成功", nil)
	return
}

func (ac *AssetController) PostDeleteAsset(c *gin.Context) {
	common.Response(c, http.StatusBadRequest, "删除资产", nil)
	return
}

func (ac *AssetController) PostUpdateAsset(c *gin.Context) {
	common.Response(c, http.StatusBadRequest, "更新资产", nil)
	return
}

func (ac *AssetController) PostFetchAsset(c *gin.Context) {
	common.Response(c, http.StatusBadRequest, "获取资产", nil)
	return
}

func (ac *AssetController) PostFetchAssetIpDetail(c *gin.Context) {
	common.Response(c, http.StatusBadRequest, "获取资产 IP 详情", nil)
	return
}

func (ac *AssetController) PostFetchAssetDomainDetail(c *gin.Context) {
	common.Response(c, http.StatusBadRequest, "获取资产域名详情", nil)
	return
}
