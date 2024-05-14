package model

import (
	"fmt"
	"github.com/gofrs/uuid/v5"
	"gorm.io/gorm"
	"time"
)

type SubDomainResult struct {
	UUID      uuid.UUID `json:"uuid" gorm:"primarykey;index;not null;comment:唯一标识符;"`
	TaskUUID  uuid.UUID `json:"taskUUID" gorm:"index;not null;comment:所属任务标识符"`
	Task      Task      `json:"task" gorm:"foreignKey:TaskUUID;references:UUID;comment:任务信息"`
	SubDomain string    `json:"subDomain" gorm:"index;not null;comment:子域名"`
	Title     string    `json:"title" gorm:"comment:网站标题"`
	Cname     string    `json:"cname" gorm:"comment:CNAME 解析"`
	Ips       string    `json:"ips" gorm:"comment:IP 地址"`
	Code      int       `json:"code" gorm:"comment:响应码"`
	Notes     string    `json:"notes" gorm:"comment:备注"`
	CreatedAt time.Time `json:"createdAt" gorm:"autoCreateTime;comment:创建时间"`
}

func (*SubDomainResult) TableName() string {
	return "sys_subdomain_results"
}

func (s *SubDomainResult) InsertData(db *gorm.DB) error {
	if s.TaskUUID != uuid.Nil {
		if err := db.Model(s).Where("task_uuid = ? AND sub_domain = ?", s.TaskUUID, s.SubDomain).FirstOrCreate(s).Error; err != nil {
			return fmt.Errorf("插入或查找端口扫描结果失败: %w", err)
		}
	}
	return nil
}
