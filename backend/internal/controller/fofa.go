package controller

import (
	"backend/internal/global"
	"backend/internal/model/common"
	"backend/internal/model/request"
	"backend/internal/service"
	"encoding/csv"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
	"time"
)

type FofaController struct {
	JWTRequired bool
}

func (fc *FofaController) PostFofaSearch(c *gin.Context) {
	var req request.FofaSearchRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		global.Logger.Error("PostFofaSearch 参数解析错误: ", zap.Error(err))
		common.ResponseOk(c, http.StatusBadRequest, "参数解析错误", nil)
		return
	}

	results, err := service.FofaServiceApp.FofaSearch(global.Config.Fofa, req)
	if err != nil {
		global.Logger.Error("PostFofaSearch 执行错误: ", zap.Error(err))
		common.ResponseOk(c, http.StatusInternalServerError, fmt.Sprintf("执行错误: %v", err), nil)
		return
	}

	common.ResponseOk(c, http.StatusOK, "查询数据成功", results)
	return
}

func (fc *FofaController) PostExportData(c *gin.Context) {
	var req request.FofaSearchRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		global.Logger.Error("PostExportData 参数解析错误: ", zap.Error(err))
		common.ResponseOk(c, http.StatusBadRequest, "参数解析错误", nil)
		return
	}

	results, err := service.FofaServiceApp.FofaSearch(global.Config.Fofa, req)
	if err != nil {
		global.Logger.Error("PostFofaSearch 执行错误: ", zap.Error(err))
		common.ResponseOk(c, http.StatusInternalServerError, fmt.Sprintf("执行错误: %v", err), nil)
		return
	}

	// 设置 CSV 文件名
	filename := "fofa_" + time.Now().Format("20060102") + ".csv"
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
	writer.Write([]string{"URL 链接", "标题", "IP 地址", "端口", "协议", "区域", "备案信息"})

	// 遍历数据并写入 CSV
	for _, item := range results.Results {
		if err := writer.Write([]string{item.URL, item.Title, item.IP, item.Port, item.Protocol, fmt.Sprintf("%v/%v/%v", item.Country, item.Region, item.City), item.ICP}); err != nil {
			global.Logger.Error("PostExportData 创建 CSV 文件失败: ", zap.Error(err))
			common.ResponseOk(c, http.StatusBadRequest, "创建 CSV 文件失败", nil)
			return
		}
	}
	return
}
