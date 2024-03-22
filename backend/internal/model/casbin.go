package model

import (
	"fmt"
	"gorm.io/gorm"
)

type CasbinRule struct {
	ID    uint   `gorm:"primaryKey;autoIncrement"`
	Ptype string `gorm:"size:100"`
	V0    string `gorm:"size:100"`
	V1    string `gorm:"size:100"`
	V2    string `gorm:"size:100"`
	V3    string `gorm:"size:100"`
	V4    string `gorm:"size:100"`
	V5    string `gorm:"size:100"`
}

func (*CasbinRule) TableName() string {
	return "casbin_rule"
}
func (c *CasbinRule) InsertData(db *gorm.DB) error {
	err := db.Where(&CasbinRule{Ptype: c.Ptype, V0: c.V0, V1: c.V1, V2: c.V2}).FirstOrCreate(c).Error // 使用 FirstOrCreate 来避免重复创建
	if err != nil {
		return fmt.Errorf("插入或查找路由失败: %w", err)
	}
	return nil
}
