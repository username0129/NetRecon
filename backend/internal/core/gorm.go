package core

import (
	"backend/internal/global"
	"gorm.io/gorm"
)

// InitializeDB
//
//	@Description: 使用 gorm 获取数据库链接
//	@return *gorm.DB
func InitializeDB() *gorm.DB {
	switch global.Config.System.DBType {
	case "mysql":
		return InitializeMysql()
	default:
		return InitializeMysql()
	} // 根据系统使用的数据库类型返回对应的数据库链接。默认使用 MySQL
}
