package model

import (
	"fmt"
	"gorm.io/gorm"
)

type Route struct {
	gorm.Model
	ParentId    uint        `json:"parentId" gorm:"comment:父路由 ID"`    // 父路由ID
	Name        string      `json:"name" gorm:"comment:路由 name"`       // 路由name
	Path        string      `json:"path" gorm:"comment:路由 path"`       // 路由path
	Meta        Meta        `json:"meta" gorm:"embedded;comment:附加属性"` // 附加属性
	Component   string      `json:"component" gorm:"comment:对应前端文件路径"` // 前端文件路径
	Authorities []Authority `json:"authorities" gorm:"many2many:sys_authority_route;"`
	Children    []Route     `json:"children" gorm:"-"`
}

type Meta struct {
	Title     string `json:"title" gorm:"comment:路由名"`                    // 路由名
	Icon      string `json:"icon" gorm:"comment路由图标"`                     // 路由图标
	KeepAlive bool   `json:"keepAlive" gorm:"default:false;comment:是否缓存"` // 是否缓存
	Hidden    bool   `json:"hidden" gorm:"default:false;comment:是否在列表隐藏"` // 是否在列表隐藏
}

func (r *Route) TableName() string {
	return "sys_routes"
}

func (r *Route) InsertData(db *gorm.DB) error {
	// 根据角色名查询数据库中已有的角色
	var uniqueAuthorities []Authority
	for _, auth := range r.Authorities {
		var tempAuth Authority
		err := db.Where("authority_name = ?", auth.AuthorityName).FirstOrCreate(&tempAuth).Error
		if err != nil {
			return fmt.Errorf("处理Authority失败: %w", err)
		}
		uniqueAuthorities = append(uniqueAuthorities, tempAuth)
	}
	r.Authorities = uniqueAuthorities

	err := db.Where(&Route{Name: r.Name}).FirstOrCreate(r).Error // 使用 FirstOrCreate 避免重复创建
	if err != nil {
		return fmt.Errorf("插入或查找路由失败: %w", err)
	}

	// 更新 many2many 关系
	if len(r.Authorities) > 0 {
		for _, authority := range r.Authorities {
			err = db.Model(r).Association("Authorities").Append(&authority)
			if err != nil {
				return fmt.Errorf("更新角色路由关系失败: %w", err)
			}
		}
	}
	return nil
}
