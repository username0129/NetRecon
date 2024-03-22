package core

import (
	"backend/internal/global"
	"backend/internal/util"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func InitializeMysql() *gorm.DB {
	cfg := global.Config.Mysql
	if err := cfg.Check(); err != nil {
		global.Logger.Error(err.Error())
		return nil
	} // 未经过初始化，返回空连接。

	mysqlConfig := mysql.Config{
		DSN:                       cfg.GetDSN(), // DSN data source name
		DefaultStringSize:         256,          // string 类型字段的默认长度
		DisableDatetimePrecision:  true,         // 禁用 datetime 精度，MySQL 5.6 之前的数据库不支持
		DontSupportRenameIndex:    true,         // 重命名索引时采用删除并新建的方式，MySQL 5.7 之前的数据库和 MariaDB 不支持重命名索引
		DontSupportRenameColumn:   true,         // 用 `change` 重命名列，MySQL 8 之前的数据库和 MariaDB 不支持重命名列
		SkipInitializeWithVersion: false,        // 根据当前 MySQL 版本自动配置
	}

	if db, err := gorm.Open(mysql.New(mysqlConfig), util.GetGormConfig()); err != nil {
		return nil // 建立连接失败
	} else {
		sqlDB, _ := db.DB()                           // 获取通用数据库对象 sql.DB。
		sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)       // 设置连接池中空闲连接的最大数量。
		sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)       // 设置打开数据库连接的最大数量。
		sqlDB.SetConnMaxLifetime(cfg.ConnMaxLifetime) // 设置空闲连接的存活时间。
		return db                                     // 建立连接成功。
	}
}
