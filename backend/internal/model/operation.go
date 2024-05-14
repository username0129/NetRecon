package model

import (
	"errors"
	"github.com/gofrs/uuid/v5"
	"gorm.io/gorm"
	"time"
)

// OperationRecord 用于记录用户操作记录
type OperationRecord struct {
	UUID      uuid.UUID `json:"uuid" gorm:"primarykey;index;not null;comment:记录标号"`           // 记录编号
	UserUUID  uuid.UUID `json:"userUUID" gorm:"index;not null;comment:用户编号"`                  // 用户编号
	User      User      `json:"user" gorm:"foreignKey:UserUUID;references:UUID;comment:用户信息"` // 用户
	IP        string    `json:"ip" gorm:"comment:客户端 IP"`                                     // 客户端 IP 地址
	Method    string    `json:"method" gorm:"comment:请求方法"`                                   // 请求方法
	Duration  string    `json:"duration" gorm:"comment:处理时间"`                                 // 处理时间
	Path      string    `json:"path" gorm:"type:text;comment:请求路径"`                           // 请求路径
	Code      string    `json:"code" gorm:"comment:请求响应状态"`                                   // 请求状态
	Agent     string    `json:"agent" gorm:"type:text;comment:浏览器代理"`                         // 代理
	Body      string    `json:"body"  gorm:"type:text;comment:请求体"`                           // 请求Body
	Resp      string    `json:"resp" gorm:"type:text;comment:响应体"`                            // 响应Body
	CreatedAt time.Time `json:"createdAt" gorm:"autoCreateTime;comment:创建时间"`                 // 创建时间
}

func (*OperationRecord) TableName() string {
	return "sys_operation_records"
}

func (o *OperationRecord) InsertData(db *gorm.DB) error {
	if o.UUID != uuid.Nil {
		if err := db.Model(&OperationRecord{}).Where("uuid = ?", o.UUID).FirstOrCreate(o).Error; err != nil {
			return errors.New("插入数据失败")
		}
	}
	return nil
}
