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

type SubDomainController struct {
	JWTRequired bool
}

func (sc *SubDomainController) PostBruteSubdomains(c *gin.Context) {
	var req request.SubDomainRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		global.Logger.Error("PostBruteSubdomains 参数解析错误: ", zap.Error(err))
		common.ResponseOk(c, http.StatusBadRequest, "参数解析错误", nil)
		return
	}
	err := service.SubDomainServiceApp.BruteSubdomains(req, util.GetUUID(c), "BruteSubdomain")
	if err != nil {
		common.ResponseOk(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}
	common.ResponseOk(c, http.StatusOK, "任务提交成功", nil)
	return
}

func (sc *SubDomainController) PostExportData(c *gin.Context) {
	var req request.DeleteSubDomainResultRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		global.Logger.Error("PostExportData 参数解析错误: ", zap.Error(err))
		common.ResponseOk(c, http.StatusBadRequest, "参数解析错误", nil)
		return
	}

	result, err := service.SubDomainServiceApp.FetchAllResult(global.DB, req.UUID)
	if err != nil {
		global.Logger.Error("PostExportData 查询数据失败: ", zap.Error(err))
		common.ResponseOk(c, http.StatusBadRequest, "查询数据失败", nil)
		return
	}

	// 设置 CSV 文件名
	filename := "subdomain_" + time.Now().Format("20060102") + ".csv"
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
	writer.Write([]string{"子域名", "响应码", "标题", "CNAME 解析", "IP 地址解析", "备注"})

	// 遍历数据并写入 CSV
	for _, subDomain := range result {
		if err := writer.Write([]string{subDomain.SubDomain, strconv.Itoa(subDomain.Code), subDomain.Title, subDomain.Cname, subDomain.Ips, subDomain.Notes}); err != nil {
			global.Logger.Error("PostExportData 创建 CSV 文件失败: ", zap.Error(err))
			common.ResponseOk(c, http.StatusBadRequest, "创建 CSV 文件失败", nil)
			return
		}
	}
	return
}

func (sc *SubDomainController) PostFetchResult(c *gin.Context) {
	var req request.FetchSubDomainResultRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		global.Logger.Error("PostFetchResult 参数解析错误: ", zap.Error(err))
		common.ResponseOk(c, http.StatusBadRequest, "参数解析错误", nil)
		return
	}

	result, total, err := service.SubDomainServiceApp.FetchResult(global.DB, req.SubDomainResult, req.PageInfo, req.OrderKey, req.Desc)
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
func (sc *SubDomainController) PostDeleteResult(c *gin.Context) {
	var req request.DeleteSubDomainResultRequest
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
