package model

import (
	"fmt"
	"github.com/gofrs/uuid/v5"
	"gorm.io/gorm"
)

type File struct {
	gorm.Model
	UUID uuid.UUID `json:"uuid" gorm:"index;comment:文件 UUID"`
	Name string    `json:"filename" gorm:"comment:文件名"` // 文件名
	Url  string    `json:"url" gorm:"comment:文件地址"`     // 文件地址
}

func (*File) TableName() string {
	return "sys_files"
}

func (f *File) InsertData(db *gorm.DB) (err error) {
	if f.UUID != uuid.Nil {
		if err := db.Model(f).Where("uuid = ? ", f.UUID).FirstOrCreate(f).Error; err != nil {
			return fmt.Errorf("插入或查找文件失败: %w", err)
		}
	}
	return nil
}
