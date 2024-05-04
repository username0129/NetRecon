package model

import (
	"fmt"
	"github.com/gofrs/uuid/v5"
	"gorm.io/gorm"
)

type SubDomainResult struct {
	gorm.Model
	UUID      uuid.UUID `json:"uuid" gorm:"index;comment:'UUID';"`
	TaskUUID  uuid.UUID `json:"taskUUID" gorm:"comment:'所属任务 UUID'"`
	Task      Task      `json:"task" gorm:"foreignKey:TaskUUID;references:UUID"` // 创建者详细信息
	SubDomain string    `json:"subDomain" gorm:"comment:'子域名'"`
	Title     string    `json:"title" gorm:"comment:'网站标题'"`
	Cname     string    `json:"cname" gorm:"comment:'CNAME 解析'"`
	Ips       string    `json:"ips" gorm:"comment:'A 解析'"`
	Code      int       `json:"code" gorm:"comment:'响应码'"`
	Notes     string    `json:"notes" gorm:"comment:'备注'"`
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
