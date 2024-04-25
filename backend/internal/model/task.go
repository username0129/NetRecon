package model

import (
	"backend/internal/global"
	"context"
	"github.com/gofrs/uuid/v5"
	"gorm.io/gorm"
)

type Task struct {
	gorm.Model
	UUID     uuid.UUID          `json:"uuid" gorm:"comment:任务 UUID;"`
	UserUUID uuid.UUID          `json:"userUUID"`
	Type     string             `json:"type" gorm:"comment:任务类型;"`                  // 例如端口扫描、子域名爆破等
	Status   string             `json:"status" gorm:"default:Running;comment:任务状态"` // 任务状态：进行中、取消、完成
	Cancel   context.CancelFunc `gorm:"-"`                                          // 不存储在数据库中，仅运行时使用
	Ctx      context.Context    `gorm:"-"`                                          // 同上
}

// CreateOrUpdate 保存或更新任务状态到数据库
func (t *Task) CreateOrUpdate() error {
	return global.DB.Save(t).Error
}

// UpdateStatus 更新任务的状态
func (t *Task) UpdateStatus(status string) error {
	t.Status = status
	return t.CreateOrUpdate()
}
