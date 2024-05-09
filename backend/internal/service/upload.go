package service

import (
	"backend/internal/global"
	"backend/internal/model"
	"backend/internal/model/request"
	"backend/internal/util/upload"
	"errors"
	"fmt"
	"github.com/gofrs/uuid/v5"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"mime/multipart"
	"strings"
)

type UploadService struct{}

var (
	UploadServiceApp = new(UploadService)
)

// extractFilenameWithoutExt 从完整文件名中提取没有扩展名的部分
func extractFilenameWithoutExt(filename string) string {
	// 查找最后一个点的位置
	dotIndex := strings.LastIndex(filename, ".")
	if dotIndex == -1 {
		// 如果没有找到点，说明文件名中没有扩展名
		return filename
	}
	// 返回点之前的部分
	return filename[:dotIndex]
}

func (fs *UploadService) UploadFile(header *multipart.FileHeader) (string, error) {
	var f model.File
	uploadTool := &upload.Local{}
	// 执行文件上传
	path, err := uploadTool.UploadFile(header)
	if err != nil {
		return "", err
	}
	// 保存文件
	f = model.File{
		UUID: uuid.Must(uuid.NewV4()),
		Url:  path,
		Name: extractFilenameWithoutExt(header.Filename),
	}
	if err := f.InsertData(global.DB); err != nil {
		return "", errors.New("保存文件信息失败")
	}
	return path, nil
}

func (fs *UploadService) DeleteFile(db *gorm.DB, fileUUID uuid.UUID) error {
	var f model.File
	uploadTool := &upload.Local{}
	if err := db.Model(&model.File{}).Where("uuid = ?", fileUUID).First(&f).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("未查询到文件信息")
		}
		global.Logger.Error("查询文件信息失败: ", zap.Error(err))
		return errors.New("查询文件信息失败")
	}

	if err := uploadTool.DeleteFile(f.Url); err != nil {
		return err
	}

	if err := db.Model(&model.File{}).Where("uuid = ?", fileUUID).Delete(&model.File{}).Error; err != nil {
		global.Logger.Error("删除文件信息失败: ", zap.Error(err))
		return errors.New("删除文件信息失败")
	}
	return nil
}

func (fs *UploadService) FetchFiles(cdb *gorm.DB, req model.File, info request.PageInfo, order string, desc bool) ([]model.File, int64, error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	db := cdb.Model(&model.File{})

	// 条件查询
	if req.Name != "" {
		db = db.Where("name LIKE ?", "%"+req.Name+"%")
	}

	// 获取满足条件的条目总数
	var total int64
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if total == 0 {
		return nil, 0, nil
	}

	// 根据有效列表进行排序处理
	orderStr := "created_at desc" // 默认排序
	if order != "" {
		allowedOrders := map[string]bool{
			"created_at": true,
		}
		if _, ok := allowedOrders[order]; !ok {
			return nil, 0, fmt.Errorf("非法的排序字段: %v", order)
		}
		orderStr = order
		if desc {
			orderStr += " desc"
		}
	}

	// 查询数据
	var resultList []model.File
	if err := db.Limit(limit).Offset(offset).Order(orderStr).Find(&resultList).Error; err != nil {
		return nil, 0, err
	}

	return resultList, total, nil
}
