package model

import (
	"fmt"
	"github.com/gofrs/uuid/v5"
	"gorm.io/gorm"
	"time"
)

type Asset struct {
	UUID        uuid.UUID `json:"uuid" gorm:"primarykey;index;not null;comment:资产的唯一标识符 UUID"`
	CreatorUUID uuid.UUID `json:"creatorUUID" gorm:"index;not null;comment:创建者 UUID;"`
	Creator     User      `json:"creator" gorm:"foreignKey:CreatorUUID;references:UUID;comment:创建者详细信息"`
	Title       string    `json:"title" gorm:"not null;comment:资产标题"`
	Domains     string    `json:"domains" gorm:"comment:资产关联的域名"`
	IPs         string    `json:"ips" gorm:"comment:资产关联的 IP"`
	CreatedAt   time.Time `json:"createdAt" gorm:"autoCreateTime;comment:创建时间"`
}

func (*Asset) TableName() string {
	return "sys_assets"
}

func (na *Asset) InsertData(db *gorm.DB) error {
	if na.UUID != uuid.Nil {
		if err := db.Model(&Asset{}).Where("uuid = ?", na.UUID).FirstOrCreate(na).Error; err != nil {
			return fmt.Errorf("插入或查找资产结果失败: %w", err)
		}
	}
	return nil
}
