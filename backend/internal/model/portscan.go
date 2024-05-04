package model

import (
	"fmt"
	"github.com/gofrs/uuid/v5"
	"gorm.io/gorm"
)

type PortScanResult struct {
	gorm.Model
	UUID     uuid.UUID `json:"uuid" gorm:"index;comment:'UUID';"`
	TaskUUID uuid.UUID `json:"taskUUID" gorm:"index;comment:'所属任务 UUID';"`      // 外键
	Task     Task      `json:"task" gorm:"foreignKey:TaskUUID;references:UUID"` // 创建者详细信息
	IP       string    `json:"ip" gorm:"comment:'目标 IP';"`                      // IP 地址
	Port     int       `json:"port" gorm:"comment:'端口';"`                       // 端口
	Service  string    `json:"service" gorm:"comment:'服务';"`                    // 端口是否开启
	Open     bool      `json:"open" gorm:"comment:'端口开启状态';"`                   // 端口是否开启
}

func (*PortScanResult) TableName() string {
	return "sys_port_scan_results"
}

func (p *PortScanResult) InsertData(db *gorm.DB) error {
	if p.TaskUUID != uuid.Nil && p.Open == true {
		if err := db.Model(p).Where("task_uuid = ? AND ip = ? AND port = ?", p.TaskUUID, p.IP, p.Port).FirstOrCreate(p).Error; err != nil {
			return fmt.Errorf("插入或查找端口扫描结果失败: %w", err)
		}
	}
	return nil
}
