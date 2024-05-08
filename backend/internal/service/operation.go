package service

import (
	"backend/internal/global"
	"backend/internal/model"
	"backend/internal/model/request"
	"errors"
	"fmt"
	"github.com/gofrs/uuid/v5"
	"gorm.io/gorm"
)

type OperationService struct {
}

var (
	OperationServiceApp = new(OperationService)
)

func (os *OperationService) FetchResult(cdb *gorm.DB, result model.OperationRecord, info request.PageInfo, order string, desc bool) ([]model.OperationRecord, int64, error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	db := cdb.Model(&model.OperationRecord{})
	// 条件查询
	if result.UserUUID != uuid.Nil {
		db = db.Where("user_uuid = ?", result.UserUUID.String())
	}
	if result.Method != "" {
		db = db.Where("method = ?", result.Method)
	}
	if result.IP != "" {
		db = db.Where("ip = ?", result.IP)
	}
	if result.Code != "" {
		db = db.Where("code = ?", result.Code)
	}
	if result.Path != "" {
		db = db.Where("path LIKE ?", "%"+result.Path+"%")
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
	orderStr := "user_uuid desc" // 默认排序
	if order != "" {
		allowedOrders := map[string]bool{
			"user_uuid":  true,
			"method":     true,
			"code":       true,
			"path":       true,
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
	var resultList []model.OperationRecord
	if err := db.Limit(limit).Offset(offset).Order(orderStr).Preload("User").Find(&resultList).Error; err != nil {
		return nil, 0, err
	}

	return resultList, total, nil
}

// DeleteResults DeleteResult 删除端口扫描结果
func (os *OperationService) DeleteResults(uuids []uuid.UUID) error {
	// 开启事务
	tx := global.DB.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	// 直接尝试删除记录
	result := tx.Where("uuid in (?)", uuids).Delete(&model.OperationRecord{})
	if result.Error != nil {
		tx.Rollback() // 如果删除操作出错，回滚事务
		return result.Error
	}

	if result.RowsAffected == 0 {
		tx.Rollback() // 如果没有删除任何记录，回滚事务
		return errors.New("没有找到记录")
	}

	// 提交事务
	return tx.Commit().Error
}
