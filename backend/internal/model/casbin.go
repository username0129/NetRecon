package model

import (
	"fmt"
	"gorm.io/gorm"
)

type CasbinRule struct {
	ID    uint   `json:"id" gorm:"primaryKey;autoIncrement;comment:唯一标识符"`
	Ptype string `json:"ptype" gorm:"not null;comment:策略类型"`
	V0    string `json:"v0" gorm:"comment:角色"`
	V1    string `json:"v1" gorm:"comment:请求路径"`
	V2    string `json:"v2" gorm:"comment:请求方法"`
	V3    string `json:"v3" gorm:"comment:策略字段3"`
	V4    string `json:"v4" gorm:"comment:策略字段4"`
	V5    string `json:"v5" gorm:"comment:策略字段5"`
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
