package model

import (
	"fmt"
	"github.com/gofrs/uuid/v5"
	"gorm.io/gorm"
)

type PortScanResult struct {
	gorm.Model
	TaskUUID uuid.UUID `json:"taskUUID" gorm:"index;comment:所属任务 UUID;"` // 外键
	IP       string    `json:"ip" gorm:"comment:目标 IP;"`                 // IP 地址
	Port     uint16    `json:"port" gorm:"comment:端口;"`                  // 端口l
	Protocol string    `json:"protocol" gorm:"comment:协议;"`              // 端口协议
	Service  string    `json:"service" gorm:"comment:服务;"`               // 端口是否开启
	Open     bool      `json:"open" gorm:"comment:端口开启状态;"`              // 端口是否开启
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
