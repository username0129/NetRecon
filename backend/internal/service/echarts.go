package service

import (
	"backend/internal/global"
	"backend/internal/model"
	"backend/internal/model/response"
	"errors"
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

func (es *EchartsService) FetchTasksCount(db *gorm.DB, userUUID uuid.UUID) ([]response.EchartsResponse, error) {
	endDate := time.Now()
	startDate := endDate.AddDate(0, 0, -4)

	// 生成完整的日期范围
	dateRange := GenerateDateRange(startDate, endDate)

	// 定义查询结果的结构
	var results []struct {
		Date  time.Time
		Count int
	}

	err := db.Model(model.Task{}).
		Select("DATE(created_at) as date, COUNT(*) as count").
		Where("creator_uuid = ? AND created_at >= ?", userUUID, startDate).
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
	mergedResults := make([]response.EchartsResponse, len(dateRange))
	for i, date := range dateRange {
		dateStr := date.Format("2006-01-02")
		count, exists := resultMap[dateStr]
		if !exists {
			count = 0 // 如果没有数据，设置 Count 为 0
		}
		mergedResults[i] = response.EchartsResponse{Date: dateStr, Count: count}
	}

	return mergedResults, nil
}
