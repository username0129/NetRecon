package controller

import (
	"backend/internal/global"
	"backend/internal/model/common"
	"backend/internal/model/request"
	"backend/internal/model/response"
	"backend/internal/service"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
)

type UploadController struct {
	JWTRequired bool
}

func (fc *UploadController) PostUploadFile(c *gin.Context) {
	_, header, err := c.Request.FormFile("file")
	if err != nil {
		common.ResponseOk(c, http.StatusBadRequest, "参数解析失败", nil)
		return
	}
	path, err := service.UploadServiceApp.UploadFile(header)
	if err != nil {
		common.ResponseOk(c, http.StatusBadRequest, fmt.Sprintf("文件上传失败: %v", err), nil)
		return
	}
	common.ResponseOk(c, http.StatusOK, "文件上传成功", gin.H{"url": path})
}

func (fc *UploadController) PostDeleteFile(c *gin.Context) {
	var file request.UUIDRequest

	if err := c.ShouldBindJSON(&file); err != nil {
		global.Logger.Error("DeleteFile 参数解析错误: ", zap.Error(err))
		common.ResponseOk(c, http.StatusBadRequest, "参数解析错误", nil)
		return
	}

	if err := service.UploadServiceApp.DeleteFile(global.DB, file.UUID); err != nil {
		common.ResponseOk(c, http.StatusInternalServerError, fmt.Sprintf("文件删除失败: %v", err), nil)
		return
	}

	common.ResponseOk(c, http.StatusOK, "文件删除成功", nil)
}

// PostFetchFiles 分页查询文件信息
func (fc *UploadController) PostFetchFiles(c *gin.Context) {
	var req request.FetchFilesRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		global.Logger.Error("PostFetchFiles 参数解析错误: ", zap.Error(err))
		common.ResponseOk(c, http.StatusBadRequest, "参数解析错误", nil)
		return
	}

	result, total, err := service.UploadServiceApp.FetchFiles(global.DB, req.File, req.PageInfo, req.OrderKey, req.Desc)
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
