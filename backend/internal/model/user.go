package model

import (
	"fmt"
	"github.com/gofrs/uuid/v5"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	UUID        uuid.UUID `json:"uuid" gorm:"uniqueIndex;comment:用户 UUID;"`
	Username    string    `json:"username" gorm:"index;comment:用户登录名;"`
	Password    string    `json:"-" gorm:"comment:用户登录密码;"`
	Nickname    string    `json:"nickname" gorm:"comment:用户昵称;"`
	AuthorityId string    `json:"authorityId" gorm:"default:1;comment:用户身份 ID;"`
	Mail        string    `json:"mail" gorm:"comment:邮箱;"`
	Enable      string    `json:"enable" gorm:"default:1;comment:用户状态 1 => 正常 2 => 冻结;"`
}

func (*User) TableName() string {
	return "sys_users"
}

func (u *User) InsertData(db *gorm.DB) error {
	if u.UUID == uuid.Nil {
		u.UUID = uuid.Must(uuid.NewV4()) // 确保 UUID 被正确设置
	}
	if err := db.Model(u).Where("username = ?", u.Username).FirstOrCreate(u).Error; err != nil {
		return fmt.Errorf("插入或查找用户失败: %w", err)
	}
	return nil
}
