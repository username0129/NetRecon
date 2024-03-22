package db

import (
	"backend/internal/config"
	"backend/internal/db/dbutil"
	"backend/internal/global"
	"backend/internal/model"
	"backend/internal/util"
	"context"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type MySQLInitializer struct {
	CommonOps dbutil.CommonDBOperations
}

func NewMySQLInitializer() *MySQLInitializer {
	return &MySQLInitializer{}
}

func (mi *MySQLInitializer) CreateDatabase(c context.Context, req model.InitRequest) (context.Context, error) {
	cfg := req.ToMySQLConfig()
	c = context.WithValue(c, "config", cfg)
	if err := cfg.Check(); err != nil {
		return nil, err
	} // 配置检查失败

	dsn := cfg.GetEmptyDSN()
	createSql := fmt.Sprintf("CREATE DATABASE IF NOT EXISTS `%s` DEFAULT CHARACTER SET utf8mb4 DEFAULT COLLATE utf8mb4_general_ci;", cfg.DBName)
	if err := dbutil.ExecuteSQL(dsn, "mysql", createSql); err != nil {
		return nil, err
	} // 创建数据库

	mysqlConfig := mysql.Config{
		DSN:                       cfg.GetDSN(), // DSN data source name
		DefaultStringSize:         256,          // string 类型字段的默认长度
		DisableDatetimePrecision:  true,         // 禁用 datetime 精度，MySQL 5.6 之前的数据库不支持
		DontSupportRenameIndex:    true,         // 重命名索引时采用删除并新建的方式，MySQL 5.7 之前的数据库和 MariaDB 不支持重命名索引
		DontSupportRenameColumn:   true,         // 用 `change` 重命名列，MySQL 8 之前的数据库和 MariaDB 不支持重命名列
		SkipInitializeWithVersion: false,        // 根据当前 MySQL 版本自动配置
	}

	if db, err := gorm.Open(mysql.New(mysqlConfig), util.GetGormConfig()); err != nil {
		return c, err
	} else {
		sqlDB, _ := db.DB()                           // 获取通用数据库对象 sql.DB。
		sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)       // 设置连接池中空闲连接的最大数量。
		sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)       // 设置打开数据库连接的最大数量。
		sqlDB.SetConnMaxLifetime(cfg.ConnMaxLifetime) // 设置空闲连接的存活时间。
		c = context.WithValue(c, "db", db)
		return c, nil
	}
}

// CreateTable 使用CommonOps中的实现
func (mi *MySQLInitializer) CreateTable() error {
	return mi.CommonOps.CreateTable()
}

func (mi *MySQLInitializer) InsertData() error {
	return mi.CommonOps.InsertData()
}

func (mi *MySQLInitializer) WriteConfig(c context.Context) error {
	mysqlConfig := c.Value("config").(config.Mysql)

	global.Config.System.DBType = "mysql"
	global.Config.Mysql = mysqlConfig

	maps := util.StructToMap(global.Config)
	for k, v := range maps {
		global.Viper.Set(k, v)
	}
	return global.Viper.WriteConfig()
}
