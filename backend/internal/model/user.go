package model

import (
	"fmt"
	"github.com/gofrs/uuid/v5"
	"gorm.io/gorm"
	"time"
)

type User struct {
	UUID        uuid.UUID `json:"uuid" gorm:"primarykey;index;not null;comment:唯一标识符"`
	Username    string    `json:"username" gorm:"index;comment:用户登录名;"`
	Password    string    `json:"-" gorm:"comment:用户登录密码;"`
	Nickname    string    `json:"nickname" gorm:"comment:用户昵称;"`
	Avatar      string    `json:"avatar" gorm:"comment:用户头像路径"`
	AuthorityId string    `json:"authorityId" gorm:"default:1;comment:用户角色编号;"`
	Mail        string    `json:"mail" gorm:"comment:邮箱;"`
	Enable      string    `json:"enable" gorm:"default:1;comment:用户状态 1 => 正常 2 => 冻结;"`
	CreatedAt   time.Time `json:"createdAt" gorm:"autoCreateTime;comment:创建时间"`
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
