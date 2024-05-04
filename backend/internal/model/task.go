package model

import (
	"backend/internal/global"
	"context"
	"fmt"
	"github.com/gofrs/uuid/v5"
	"gorm.io/gorm"
)

type Task struct {
	gorm.Model
	UUID        uuid.UUID          `json:"uuid" gorm:"uniqueIndex;comment:任务 UUID;"`
	CreatorUUID uuid.UUID          `json:"creatorUUID" gorm:"index;comment:创建者 UUID;"`            // 创建者 UUID
	Creator     User               `json:"creator" gorm:"foreignKey:CreatorUUID;references:UUID"` // 创建者详细信息
	Targets     string             `json:"targets" gorm:"comment:任务目标;"`
	Title       string             `json:"title" gorm:"comment:任务标题;"`
	Type        string             `json:"type" gorm:"comment:任务类型;"`                                                      // 任务类型，例如端口扫描、子域名爆破等
	DictType    string             `json:"dictType" gorm:"comment:字典类型"`                                                   // 字典类型
	Status      string             `json:"status" gorm:"default:1;comment:'任务状态，1 -> 进行中, 2 -> 已完成, 3 -> 已取消，4 -> 执行失败''"` // 任务状态：1 -> 进行中, 2 -> 已完成, 3 -> 已取消，4 -> 执行失败
	Cancel      context.CancelFunc `json:"-" gorm:"-"`                                                                     // 不存储在数据库中，仅运行时使用
	Ctx         context.Context    `json:"-" gorm:"-"`                                                                     // 同上
}

func (*Task) TableName() string {
	return "sys_tasks"
}

func (t *Task) InsertData(db *gorm.DB) (err error) {
	if t.UUID != uuid.Nil {
		if err := db.Model(t).Where("uuid = ? ", t.UUID).FirstOrCreate(t).Error; err != nil {
			return fmt.Errorf("插入或查找任务失败: %w", err)
		}
	}
	return nil
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
