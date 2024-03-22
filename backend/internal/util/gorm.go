package util

import (
	"backend/internal/global"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"os"
	"time"
)

func GetGormConfig() *gorm.Config {
	cfg := &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
	} // 禁用外键约束

	_default := logger.New(log.New(os.Stdout, "\r\n", log.LstdFlags), logger.Config{
		SlowThreshold: 200 * time.Millisecond, // 慢 SQL 查询阈值
		LogLevel:      logger.Warn,            // 日志级别
		Colorful:      true,                   // 彩色
	})

	var logMode string
	switch global.Config.System.DBType {
	case "mysql":
		logMode = global.Config.Mysql.LogMode
	default:
		logMode = global.Config.Mysql.LogMode
	}

	switch logMode {
	case "error":
		cfg.Logger = _default.LogMode(logger.Error)
	case "silent":
		cfg.Logger = _default.LogMode(logger.Silent)
	case "warn":
		cfg.Logger = _default.LogMode(logger.Warn)
	case "info":
		cfg.Logger = _default.LogMode(logger.Info)
	default:
		cfg.Logger = _default.LogMode(logger.Info)
	}
	return cfg
}
