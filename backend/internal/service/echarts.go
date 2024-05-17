package service

import (
	"backend/internal/global"
	"backend/internal/model"
	"backend/internal/model/response"
	"errors"
	"fmt"
	"github.com/gofrs/uuid/v5"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"time"
)

type EchartsService struct{}

var (
	EchartsServiceApp = new(EchartsService)
)

func GenerateDateRange(startDate, endDate time.Time) []time.Time {
	var dates []time.Time
	for d := startDate; !d.After(endDate); d = d.AddDate(0, 0, 1) {
		dates = append(dates, d)
	}
	return dates
}

func (es *EchartsService) FetchTasksCount(cdb *gorm.DB, userUUID uuid.UUID, authorityId string) ([]response.LineResponse, error) {
	// 生成完整的日期范围
	endDate := time.Now()
	startDate := endDate.AddDate(0, 0, -4)
	dateRange := GenerateDateRange(startDate, endDate)

	db := cdb.Model(&model.Task{})

	// 管理员用户可查看全部扫描任务
	if authorityId != "1" {
		db = db.Where("creator_uuid = ? AND created_at >= ?", userUUID.String(), startDate)
	} else {
		db = db.Where("created_at >= ?", startDate)
	}

	// 定义查询结果的结构
	var results []struct {
		Date  time.Time
		Count int
	}

	err := db.Select("DATE(created_at) as date, COUNT(*) as count").
		Group("date").
		Find(&results).Error
	if err != nil {
		global.Logger.Error("根据日期获取数量失败", zap.Error(err))
		return nil, errors.New("根据日期获取数量失败")
	}

	// 将查询结果转换为映射表
	resultMap := make(map[string]int)
	for _, result := range results {
		resultMap[result.Date.Format("2006-01-02")] = result.Count
	}

	// 合并日期和查询结果
	mergedResults := make([]response.LineResponse, len(dateRange))
	for i, date := range dateRange {
		dateStr := date.Format("2006-01-02")
		count, exists := resultMap[dateStr]
		if !exists {
			count = 0 // 如果没有数据，设置 Count 为 0
		}
		mergedResults[i] = response.LineResponse{Date: dateStr, Count: count}
	}

	return mergedResults, nil
}

func (es *EchartsService) FetchPortCount(db *gorm.DB, userUUID uuid.UUID, authorityId string) ([]response.PieResponse, error) {
	query := db.Model(&model.Task{})

	// 管理员用户可查看全部扫描任务
	if authorityId != "1" {
		query = query.Where("creator_uuid = ?", userUUID.String())
	}

	// 获取用户下发的所有任务
	var tasks []struct {
		UUID uuid.UUID
	}

	if err := query.Select("uuid").
		Where("type LIKE ?", "%Port%").
		Find(&tasks).Error; err != nil {
		global.Logger.Error("获取任务信息失败", zap.String("authorityId", authorityId), zap.String("userUUID", userUUID.String()), zap.Error(err))
		return nil, fmt.Errorf("获取任务信息失败: %w", err)
	}

	var results []response.PieResponse
	var finalResult []response.PieResponse

	for _, task := range tasks {
		// 获取 Port 任务相关的数据
		if err := db.Model(&model.PortScanResult{}).
			Select("ip as target,COUNT(*) as count").
			Where("task_uuid = ?", task.UUID).
			Group("target").
			Find(&results).Error; err != nil {
			global.Logger.Error("获取数量失败", zap.Error(err))
			return nil, errors.New("获取数量失败")
		}
		finalResult = append(finalResult, results...)
	}

	return finalResult, nil
}

func (es *EchartsService) FetchDomainCount(db *gorm.DB, userUUID uuid.UUID, authorityId string) ([]response.PieResponse, error) {
	query := db.Model(&model.Task{})

	// 管理员用户可查看全部扫描任务
	if authorityId != "1" {
		query = query.Where("creator_uuid = ?", userUUID.String())
	}

	// 获取用户下发的所有任务
	var tasks []struct {
		UUID uuid.UUID
	}

	if err := query.Select("uuid").
		Where("type LIKE ?", "%Domain%").
		Find(&tasks).Error; err != nil {
		global.Logger.Error("获取任务信息失败", zap.String("authorityId", authorityId), zap.String("userUUID", userUUID.String()), zap.Error(err))
		return nil, fmt.Errorf("获取任务信息失败: %w", err)
	}

	var results []response.PieResponse
	var finalResult []response.PieResponse

	for _, task := range tasks {
		// 获取 Domain 任务相关的数据
		if err := db.Model(&model.SubDomainResult{}).
			Select("domain as target,COUNT(*) as count").
			Where("task_uuid = ?", task.UUID).
			Group("target").
			Find(&results).Error; err != nil {
			global.Logger.Error("获取数量失败", zap.Error(err))
			return nil, errors.New("获取数量失败")
		}
		finalResult = append(finalResult, results...)
	}

	return finalResult, nil
}
