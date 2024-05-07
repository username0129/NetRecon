package controller

import (
	"backend/internal/global"
	"backend/internal/model/common"
	"backend/internal/model/request"
	"backend/internal/model/response"
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
		common.Response(c, http.StatusInternalServerError, fmt.Sprintf("添加资产错误：%v", err.Error()), nil)
		return
	}
	common.Response(c, http.StatusOK, "添加资产成功", nil)
	return
}

func (ac *AssetController) PostDeleteAsset(c *gin.Context) {
	var req request.DeleteAssetRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		global.Logger.Error("PostDeleteAsset 参数解析错误: ", zap.Error(err))
		common.Response(c, http.StatusBadRequest, "参数解析错误", nil)
		return
	}

	err := service.AssetServiceApp.DeleteAsset(global.DB, req.UUID, util.GetUUID(c), util.GetAuthorityId(c))
	if err != nil {
		global.Logger.Error("删除资产失败: ", zap.Error(err))
		common.Response(c, http.StatusInternalServerError, fmt.Sprintf("删除资产失败: %v", err), nil)
		return
	}
	common.Response(c, http.StatusOK, "删除资产成功", nil)
	return
}

func (ac *AssetController) PostUpdateAsset(c *gin.Context) {
	var req request.UpdateAssetRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		global.Logger.Error("PostUpdateAsset 参数解析错误: ", zap.Error(err))
		common.Response(c, http.StatusBadRequest, "参数解析错误", nil)
		return
	}

	if err := service.AssetServiceApp.UpdateAsset(global.DB, req); err != nil {
		global.Logger.Error("更新资产错误: ", zap.Error(err))
		common.Response(c, http.StatusInternalServerError, fmt.Sprintf("更新资产错误：%v", err.Error()), nil)
		return
	}
	common.Response(c, http.StatusOK, "更新资产成功", nil)
	return
}

func (ac *AssetController) PostFetchAsset(c *gin.Context) {
	var req request.FetchAssetsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		global.Logger.Error("PostFetchAsset 参数解析错误: ", zap.Error(err))
		common.Response(c, http.StatusBadRequest, "参数解析错误", nil)
		return
	}

	result, total, err := service.AssetServiceApp.FetchAssets(global.DB, req.Asset, req.PageInfo, req.OrderKey, req.Desc, util.GetUUID(c), util.GetAuthorityId(c))
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

func (ac *AssetController) PostFetchAssetIpDetail(c *gin.Context) {
	common.Response(c, http.StatusOK, "获取资产 IP 详情", nil)
	return
}

func (ac *AssetController) PostFetchAssetDomainDetail(c *gin.Context) {
	common.Response(c, http.StatusOK, "获取资产域名详情", nil)
	return
}
