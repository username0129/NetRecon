package model

import (
	"fmt"
	"gorm.io/gorm"
)

type Authority struct {
	gorm.Model
	AuthorityName string `json:"authorityName" gorm:"comment:角色名称"`
	Menus         []Menu `json:"menus" gorm:"many2many:sys_authority_menu;"`
	Users         []User `json:"users" gorm:"many2many:sys_user_authority;"`
}

func (*Authority) TableName() string {
	return "sys_authorities"
}

func (a *Authority) InsertData(db *gorm.DB) error {

	err := db.Where(&Authority{AuthorityName: a.AuthorityName}).FirstOrCreate(a).Error // 使用FirstOrCreate来避免重复创建
	if err != nil {
		return fmt.Errorf("插入或查找角色失败: %w", err)
	}

	// 更新或插入 many2many 关系
	if len(a.Menus) > 0 {
		for _, menu := range a.Menus {
			err = db.Model(a).Association("Menus").Append(&menu)
			if err != nil {
				return fmt.Errorf("更新角色菜单关系失败: %w", err)
			}
		}
	}
	return nil
}
