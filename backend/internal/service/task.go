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

type TaskService struct{}

var (
	TaskServiceApp = new(TaskService)
)

func (ts TaskService) FetchAllTasks() (tasks []model.Task, err error) {
	if err = global.DB.Find(&tasks).Error; err != nil {
		return nil, err // 返回错误信息
	}
	return tasks, nil // 返回任务列表
}

func (ts TaskService) FetchTasks(cdb *gorm.DB, result model.Task, info request.PageInfo, order string, desc bool) ([]model.Task, int64, error) {

	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	db := cdb.Model(&model.Task{})

	// 条件查询
	if result.UUID != uuid.Nil {
		db = db.Where("uuid LIKE ?", "%"+result.UUID.String()+"%")
	}
	if result.Targets != "" {
		db = db.Where("targets LIKE ?", "%"+result.Targets+"%")
	}
	if result.Title != "" {
		db = db.Where("title LIKE ?", "%"+result.Title+"%")
	}
	if result.Type != "" {
		db = db.Where("type LIKE ?", "%"+result.Type+"%")
	}
	if result.Status != "" {
		db = db.Where("status LIKE ?", "%"+result.Status+"%")
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
			"uuid":       true,
			"title":      true,
			"targets":    true,
			"type":       true,
			"status":     true,
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
	var resultList []model.Task
	if err := db.Limit(limit).Offset(offset).Order(orderStr).Find(&resultList).Error; err != nil {
		return nil, 0, err
	}

	return resultList, total, nil
}

func (ts TaskService) CancelTask(taskUUID, userUUID uuid.UUID, authorityId uint) (err error) {
	if cancel, exists := global.TaskManager[taskUUID]; exists { // 任务在管理器中存在
		var task model.Task
		if err := global.DB.Model(&task).Where("uuid = ?", taskUUID).First(&task).Error; err != nil { // 查询失败
			return err
		} else {
			if task.Status != "1" { // 任务不正在进行中
				return errors.New("任务不在运行状态")
			}
			if task.CreatorUUID != userUUID || authorityId != 1 { // 任务的发起者不是当前用户 或者 不是管理员
				return errors.New("没有权限取消任务")
			}
			// 取消任务
			cancel()
			return nil
		}
	} else {
		return errors.New("任务不存在")
	}
}
