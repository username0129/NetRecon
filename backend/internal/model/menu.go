package model

import (
	"fmt"
	"gorm.io/gorm"
)

type Menu struct {
	gorm.Model
	Hidden      bool        `json:"hidden" gorm:"comment:是否在列表隐藏"`     // 是否在列表隐藏
	ParentId    string      `json:"parentId" gorm:"comment:父菜单 ID"`    // 父菜单ID
	Title       string      `json:"title" gorm:"comment:菜单名"`          // 菜单名
	Icon        string      `json:"icon" gorm:"comment:菜单图标"`          // 菜单图标
	Name        string      `json:"name" gorm:"comment:路由 name"`       // 路由name
	Path        string      `json:"path" gorm:"comment:路由 path"`       // 路由path
	Component   string      `json:"component" gorm:"comment:对应前端文件路径"` // 对应前端文件路径
	Authorities []Authority `json:"authorities" gorm:"many2many:sys_authority_menu;"`
	Children    []Menu      `json:"children" gorm:"-"`
}

func (*Menu) TableName() string {
	return "sys_menus"
}

func (m *Menu) InsertData(db *gorm.DB) error {
	// 根据角色名查询数据库中已有的角色
	var uniqueAuthorities []Authority
	for _, auth := range m.Authorities {
		var tempAuth Authority
		err := db.Where("authority_name = ?", auth.AuthorityName).FirstOrCreate(&tempAuth).Error
		if err != nil {
			return fmt.Errorf("处理Authority失败: %w", err)
		}
		uniqueAuthorities = append(uniqueAuthorities, tempAuth)
	}
	m.Authorities = uniqueAuthorities

	err := db.Where(&Menu{Name: m.Name}).FirstOrCreate(m).Error // 使用 FirstOrCreate 避免重复创建
	if err != nil {
		return fmt.Errorf("插入或查找菜单失败: %w", err)
	}

	// 更新 many2many 关系
	if len(m.Authorities) > 0 {
		for _, authority := range m.Authorities {
			err = db.Model(m).Association("Authorities").Append(&authority)
			if err != nil {
				return fmt.Errorf("更新菜单角色关系失败: %w", err)
			}
		}
	}
	return nil
}
