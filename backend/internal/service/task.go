package service

import (
	"backend/internal/global"
	"backend/internal/model"
	"backend/internal/model/request"
	"errors"
	"fmt"
	"github.com/gofrs/uuid/v5"
	"github.com/robfig/cron/v3"
	"gorm.io/gorm"
)

type TaskService struct{}

var (
	TaskServiceApp = new(TaskService)
)

func (ts *TaskService) FetchTasks(cdb *gorm.DB, result model.Task, info request.PageInfo, order string, desc bool, userUUID uuid.UUID, authorityId string) ([]model.Task, int64, error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	db := cdb.Model(&model.Task{})

	// 管理员用户可查看全部扫描任务
	if authorityId != "1" {
		db = db.Where("creator_uuid LIKE ?", "%"+userUUID.String()+"%")
	}

	// 条件查询
	if result.UUID != uuid.Nil {
		db = db.Where("uuid LIKE ?", "%"+result.UUID.String()+"%")
	}
	if result.AssetUUID != uuid.Nil {
		db = db.Where("asset_uuid LIKE ?", "%"+result.AssetUUID.String()+"%")
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
			"dict_type":  true,
			"asset_uuid": true,
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
	if err := db.Preload("Creator").Limit(limit).Offset(offset).Order(orderStr).Find(&resultList).Error; err != nil {
		return nil, 0, err
	}

	return resultList, total, nil
}

func (ts *TaskService) FetchTaskCount(cdb *gorm.DB, taskType string, userUUID uuid.UUID, authorityId string) (int64, error) {

	db := cdb.Model(&model.Task{})

	// 管理员用户可查看全部扫描任务
	if authorityId != "1" {
		db = db.Where("creator_uuid LIKE ?", "%"+userUUID.String()+"%")
	}

	if taskType != "" {
		db = db.Where("type LIKE ?", "%"+taskType+"%")
	}

	// 获取满足条件的条目总数
	var total int64
	if err := db.Count(&total).Error; err != nil {
		return 0, err
	}

	return total, nil
}
func (ts *TaskService) CancelTask(taskUUID, userUUID uuid.UUID, authorityId string) (err error) {
	if cancel, exists := global.TaskManager[taskUUID]; exists { // 任务在管理器中存在
		var task model.Task
		// 首先从数据库中读取任务数据
		if err := global.DB.Model(&task).Where("uuid = ?", taskUUID).First(&task).Error; err != nil { // 查询失败
			return err
		}

		if task.Status != "1" { // 任务不正在进行中
			return errors.New("任务不在运行状态")
		}
		if task.CreatorUUID != userUUID && authorityId != "1" { // 任务的发起者不是当前用户 或者 不是管理员
			return errors.New("没有权限取消任务")
		}

		// 取消任务
		cancel()
		task.UpdateStatus("3") // 更新任务状态为取消

		// 如果任务属于资产监控
		if task.AssetUUID != uuid.Nil {
			// 删除定时任务
			global.CronManager.RemoveTask(cron.EntryID(task.CronID))
		}
		return nil
	} else {
		return errors.New("任务不存在")
	}
}

// DeleteTask 删除任务及其结果
func (ts *TaskService) DeleteTask(db *gorm.DB, taskUUID, userUUID uuid.UUID, authorityId string) (err error) {
	var task model.Task

	// 首先获取任务信息，确保任务存在
	if err := global.DB.Model(&model.Task{}).Where("uuid = ?", taskUUID).First(&task).Error; err != nil {
		return err // 可能是因为没有找到任务
	}

	// 检查是否有权限取消任务
	if task.Status == "1" {
		return errors.New("任务正在运行中，请等待执行完成。")
	}

	// 检查是否有权限取消任务
	if task.CreatorUUID != userUUID && authorityId != "1" {
		return errors.New("没有权限删除任务")
	}

	// 根据任务类型执行不同的删除策略
	switch task.Type {
	case "PortScan", "Corn/Port":
		// 删除 PortScan 任务相关的数据
		if err := db.Model(&model.PortScanResult{}).Where("task_uuid = ?", taskUUID).Delete(&model.PortScanResult{}).Error; err != nil {
			return fmt.Errorf("删除 PortScan 结果失败: %w", err)
		}
	case "BruteSubdomain", "Cron/Domain":
		// 删除 Subdomain 任务相关的数据
		if err := db.Model(&model.SubDomainResult{}).Where("task_uuid = ?", taskUUID).Delete(&model.SubDomainResult{}).Error; err != nil {
			return fmt.Errorf("删除 Subdomain 结果失败: %w", err)
		}
	}

	// 删除任务本身
	if result := db.Model(&model.Task{}).Where("uuid = ?", taskUUID).Delete(&model.Task{}); result.Error != nil {
		return fmt.Errorf("删除任务失败: %w", result.Error)
	}

	return nil
}

// DeleteTasks 批量任务及其结果
func (ts *TaskService) DeleteTasks(db *gorm.DB, taskUUIDs []uuid.UUID, userUUID uuid.UUID, authorityId string) (err error) {
	for _, v := range taskUUIDs {
		err := ts.DeleteTask(db, v, userUUID, authorityId)
		if !errors.Is(err, errors.New("任务正在运行中，请等待执行完成。")) && !errors.Is(err, errors.New("没有权限删除任务")) {
			return err
		}
	}
	return nil
}
