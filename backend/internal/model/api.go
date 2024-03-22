package model

import (
	"fmt"
	"gorm.io/gorm"
)

type Api struct {
	gorm.Model
	Path        string `json:"path" gorm:"comment:api 路径"`
	Description string `json:"description" gorm:"comment:api 描述"`
	Group       string `json:"group" gorm:"comment:api 组"`
	Method      string `json:"method" gorm:"default:POST;comment:请求方法"`
}

func (*Api) TableName() string {
	return "sys_apis"
}

func (a *Api) InsertData(db *gorm.DB) error {
	err := db.Where(&Api{Path: a.Path}).FirstOrCreate(a).Error // 使用FirstOrCreate来避免重复创建
	if err != nil {
		return fmt.Errorf("插入或查找路由失败: %w", err)
	}
	return nil
}
