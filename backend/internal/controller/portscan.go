package controller

import (
	"backend/internal/global"
	"backend/internal/model/common"
	"backend/internal/model/request"
	"backend/internal/model/response"
	"backend/internal/service"
	"backend/internal/util"
	"encoding/csv"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
	"strconv"
	"time"
)

type PortScanController struct {
	JWTRequired bool
}

// PostPortScan PostPort 执行端口扫描
func (pc *PortScanController) PostPortScan(c *gin.Context) {
	var portScanRequest request.PortScanRequest

	if err := c.ShouldBindJSON(&portScanRequest); err != nil {
		global.Logger.Error("PostPortScan 参数解析错误: ", zap.Error(err))
		common.ResponseOk(c, http.StatusBadRequest, "参数解析错误", nil)
		return
	}

	err := service.PortServiceApp.ExecutePortScan(portScanRequest, util.GetUUID(c), util.GetAuthorityId(c), "PortScan")
	if err != nil {
		common.ResponseOk(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}
	common.ResponseOk(c, http.StatusOK, "任务提交成功", nil)
	return
}

func (pc *PortScanController) PostFetchResult(c *gin.Context) {
	var req request.FetchPortScanResultRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		global.Logger.Error("PostFetchResult 参数解析错误: ", zap.Error(err))
		common.ResponseOk(c, http.StatusBadRequest, "参数解析错误", nil)
		return
	}

	result, total, err := service.PortServiceApp.FetchResult(global.DB, req.PortScanResult, req.PageInfo, req.OrderKey, req.Desc)
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

// PostDeleteResult 删除指定结果
func (pc *PortScanController) PostDeleteResult(c *gin.Context) {
	var req request.DeletePortScanResultRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		global.Logger.Error("PostDeleteResult 参数解析错误: ", zap.Error(err))
		common.ResponseOk(c, http.StatusBadRequest, "参数解析错误", nil)
		return
	}

	if err := service.PortServiceApp.DeleteResult(req.UUID); err != nil {
		global.Logger.Error("PostDeleteResult 运行失败: ", zap.Error(err))
		common.ResponseOk(c, http.StatusInternalServerError, "目标结果删除失败", nil)
		return
	} else {
		common.ResponseOk(c, http.StatusOK, "目标结果删除成功", nil)
		return
	}
}

// PostDeleteResults 删除指定结果
func (pc *PortScanController) PostDeleteResults(c *gin.Context) {
	var req request.DeletePortScanResultsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		global.Logger.Error("PostDeleteResult 参数解析错误: ", zap.Error(err))
		common.ResponseOk(c, http.StatusBadRequest, "参数解析错误", nil)
		return
	}

	if err := service.PortServiceApp.DeleteResults(req.UUIDS); err != nil {
		global.Logger.Error("PostDeleteResult 运行失败: ", zap.Error(err))
		common.ResponseOk(c, http.StatusInternalServerError, "目标结果删除失败", nil)
		return
	} else {
		common.ResponseOk(c, http.StatusOK, "目标结果删除成功", nil)
		return
	}
}

func (pc *PortScanController) PostExportData(c *gin.Context) {
	var req request.DeletePortScanResultRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		global.Logger.Error("PostExportData 参数解析错误: ", zap.Error(err))
		common.ResponseOk(c, http.StatusBadRequest, "参数解析错误", nil)
		return
	}

	results, err := service.PortServiceApp.FetchAllResult(global.DB, req.UUID)
	if err != nil {
		global.Logger.Error("PostExportData 查询数据失败: ", zap.Error(err))
		common.ResponseOk(c, http.StatusBadRequest, "查询数据失败", nil)
		return
	}

	// 设置 CSV 文件名
	filename := "portscan_" + time.Now().Format("20060102") + ".csv"
	// 设置响应内容类型为 CSV
	c.Writer.Header().Set("Content-type", "text/csv; charset=utf-8")
	// 设置下载的文件名
	c.Writer.Header().Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, filename))

	// 写入 UTF-8 BOM，防止中文乱码
	c.Writer.Write([]byte{0xEF, 0xBB, 0xBF})

	// 创建 csv writer 并写入数据
	writer := csv.NewWriter(c.Writer)
	defer writer.Flush()

	// 写入 CSV 头部
	writer.Write([]string{"IP 地址", "端口号", "开放服务", "链接"})

	// 遍历数据并写入 CSV
	for _, result := range results {
		if err := writer.Write([]string{result.IP, strconv.Itoa(result.Port), result.Service, fmt.Sprintf("%v://%v:%v", result.Service, result.IP, result.Port)}); err != nil {
			global.Logger.Error("PostExportData 创建 CSV 文件失败: ", zap.Error(err))
			common.ResponseOk(c, http.StatusBadRequest, "创建 CSV 文件失败", nil)
			return
		}
	}
	return
}
