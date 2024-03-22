package model

import (
	"fmt"
	"github.com/gofrs/uuid/v5"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	UUID        uuid.UUID   `json:"uuid" gorm:"comment:用户 UUID;"`
	Username    string      `json:"username" gorm:"index;comment:用户登录名;"`
	Password    string      `json:"password" gorm:"comment:用户登录密码;"`
	Nickname    string      `json:"nickname" gorm:"comment:用户昵称;"`
	Avatar      string      `json:"avatar" gorm:"comment:用户头像;"`
	AuthorityId uint        `json:"authority_id" gorm:"default:1;comment:用户身份 ID;"`
	Authorities []Authority `json:"authorities" gorm:"many2many:sys_user_authority;"`
	Email       string      `json:"email" gorm:"comment:邮箱;"`
	Enable      int         `json:"enable" gorm:"default:1;comment:用户状态 1正常 0冻结;"`
}

func (*User) TableName() string {
	return "sys_users"
}

func (u *User) InsertData(db *gorm.DB) error {
	// 根据角色名查询数据库中已有的角色
	var uniqueAuthorities []Authority
	for _, auth := range u.Authorities {
		var tempAuth Authority
		err := db.Where("authority_name = ?", auth.AuthorityName).FirstOrCreate(&tempAuth).Error
		if err != nil {
			return fmt.Errorf("处理Authority失败: %w", err)
		}
		uniqueAuthorities = append(uniqueAuthorities, tempAuth)
	}
	u.Authorities = uniqueAuthorities

	if u.UUID == uuid.Nil {
		u.UUID = uuid.Must(uuid.NewV4()) // 确保 UUID 被正确设置
	}

	err := db.Where(&User{Username: u.Username}).FirstOrCreate(u).Error // 使用 FirstOrCreate 避免重复创建
	if err != nil {
		return fmt.Errorf("插入或查找用户失败: %w", err)
	}

	// 更新 many2many 关系
	if len(u.Authorities) > 0 {
		for _, authority := range u.Authorities {
			err = db.Model(u).Association("Authorities").Append(&authority)
			if err != nil {
				return fmt.Errorf("更新用户角色关系失败: %w", err)
			}
		}
	}
	return nil
}
