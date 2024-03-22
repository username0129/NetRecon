package service

import (
	"backend/internal/db"
	"backend/internal/global"
	"backend/internal/model"
	"context"
	"gorm.io/gorm"
)

// -------------------- 数据库初始化 ----------------------------

// IDatabaseInitializer 定义数据库初始化器接口
type IDatabaseInitializer interface {
	CreateDatabase(ctx context.Context, req model.InitRequest) (context.Context, error)
	CreateTable() error
	InsertData() error
	WriteConfig(ctx context.Context) error
}

type InitService struct{}

func (is *InitService) Init(req model.InitRequest) (err error) {
	c := context.TODO()

	var dbInitializer IDatabaseInitializer
	switch req.DBType {
	case "mysql":
		dbInitializer = db.NewMySQLInitializer()
		c = context.WithValue(c, "db_type", "mysql")
	default:
		dbInitializer = db.NewMySQLInitializer()
		c = context.WithValue(c, "db_type", "mysql")
	} // 获取对应的数据库初始化工具

	if c, err = dbInitializer.CreateDatabase(c, req); err != nil {
		return err
	} // 创建数据库

	global.DB = c.Value("db").(*gorm.DB)

	if err = dbInitializer.CreateTable(); err != nil {
		return err
	}
	if err = dbInitializer.InsertData(); err != nil {
		return err
	}
	if err = dbInitializer.WriteConfig(c); err != nil {
		return err
	}
	return nil
}
