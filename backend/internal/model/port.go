package model

import (
	"errors"
	"fmt"
	"gorm.io/gorm"
)

type PortScanRequest struct {
	Targets string `json:"targets"` // 目标 IP
	Ports   string `json:"ports"`   // 目标 端口
	Timeout string `json:"timeout"` // 自定义超时时间
	Threads string `json:"threads"` // 线程数
}

type PortParseRequest struct {
	Ports string `json:"target"` // 目标 IP
}

type PortScanResult struct {
	gorm.Model
	IP       string `json:"ip" gorm:"comment:目标 IP;"`    // IP 地址
	Port     uint16 `json:"port" gorm:"comment:端口;"`     // 端口l
	Protocol string `json:"protocol" gorm:"comment:协议;"` // 端口协议
	Service  string `json:"service" gorm:"comment:服务;"`  // 端口是否开启
	Open     bool   `json:"open" gorm:"comment:是否开启;"`   // 端口是否开启
}

func (*PortScanResult) TableName() string {
	return "sys_port_scan_result"
}

func (p *PortScanResult) InsertData(db *gorm.DB) error {
	var existingResult PortScanResult
	if err := db.Where("ip = ? AND port = ?", p.IP, p.Port).First(&existingResult).Error; err != nil {
		// 如果查找失败，说明记录不存在
		if errors.Is(err, gorm.ErrRecordNotFound) {
			if p.Open {
				// 如果状态为 open，创建新记录
				if err := db.Create(p).Error; err != nil {
					return fmt.Errorf("插入新记录失败: %w", err)
				}
			}
			return nil
		} else {
			return fmt.Errorf("查找记录失败: %w", err)
		}
	}

	// 如果查找成功
	if !p.Open {
		// 如果最新扫描状态为 off，则将记录从数据库中删除
		if err := db.Delete(&existingResult).Error; err != nil {
			return fmt.Errorf("删除记录失败: %w", err)
		}
	}
	return nil
}
