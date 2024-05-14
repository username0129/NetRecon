package model

import (
	"fmt"
	"github.com/gofrs/uuid/v5"
	"gorm.io/gorm"
	"time"
)

type File struct {
	UUID      uuid.UUID `json:"uuid" gorm:"primarykey;index;not null;comment:唯一标识符"`
	Name      string    `json:"filename" gorm:"comment:文件名"` // 文件名
	Url       string    `json:"url" gorm:"comment:文件地址"`     // 文件地址
	CreatedAt time.Time `json:"createdAt" gorm:"autoCreateTime;comment:创建时间"`
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
