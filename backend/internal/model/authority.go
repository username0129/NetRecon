package model

import (
	"fmt"
	"gorm.io/gorm"
)

type Authority struct {
	gorm.Model
	AuthorityName string  `json:"authorityName" gorm:"comment:角色名称"`
	Routes        []Route `json:"routes" gorm:"many2many:sys_authority_route;"`
	Users         []User  `json:"users" gorm:"many2many:sys_user_authority;"`
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
	if len(a.Routes) > 0 {
		for _, route := range a.Routes {
			err = db.Model(&Authority{}).Association("Routes").Append(&route)
			if err != nil {
				return fmt.Errorf("更新角色路由关系失败: %w", err)
			}
		}
	}
	return nil
}
