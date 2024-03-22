package dbutil

import (
	"backend/internal/global"
	"backend/internal/model"
	"backend/internal/util"
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
)

// InitialData 用于数据库初始化的数据结构
type InitialData struct {
	TableName string
	Data      []interface{}
}

// initialDatas 定义初始化数据库时使用的数据
var initialDatas = []InitialData{
	{
		TableName: "sys_authorities",
		Data: []interface{}{
			&model.Authority{AuthorityName: "系统管理员"},
		},
	},
	{
		TableName: "casbin_role",
		Data: []interface{}{
			&model.CasbinRule{Ptype: "p", V0: "1", V1: "/api/v1/user/postuserinfo", V2: "POST"},
		},
	},
	{
		TableName: "sys_menus",
		Data: []interface{}{
			// 顶级菜单
			&model.Menu{Hidden: false, ParentId: "0", Title: "仪表盘", Icon: "odometer", Name: "dashboard", Path: "dashboard", Component: "view/dashboard/index.vue", Authorities: []model.Authority{{AuthorityName: "系统管理员"}}},
			&model.Menu{Hidden: false, ParentId: "0", Title: "管理员面板", Icon: "user", Name: "admin", Path: "admin", Component: "view/admin/index.vue", Authorities: []model.Authority{{AuthorityName: "系统管理员"}}},
			&model.Menu{Hidden: false, ParentId: "0", Title: "个人信息", Icon: "message", Name: "person", Path: "person", Component: "view/person/index.vue", Authorities: []model.Authority{{AuthorityName: "系统管理员"}}},
			// 管理员菜单
			&model.Menu{Hidden: false, ParentId: "2", Title: "角色管理", Icon: "avatar", Name: "authority", Path: "authority", Component: "view/admin/authority/index.vue", Authorities: []model.Authority{{AuthorityName: "系统管理员"}}},
			&model.Menu{Hidden: false, ParentId: "2", Title: "用户管理", Icon: "coordinate", Name: "user", Path: "user", Component: "view/admin/user/index.vue", Authorities: []model.Authority{{AuthorityName: "系统管理员"}}},
			&model.Menu{Hidden: false, ParentId: "2", Title: "操作历史", Icon: "pie-chart", Name: "operation", Path: "operation", Component: "view/admin/operation/index.vue", Authorities: []model.Authority{{AuthorityName: "系统管理员"}}},
		},
	},
	{
		TableName: "sys_apis",
		Data: []interface{}{
			// 基础 API
			&model.Api{Path: "/api/v1/base/gethealth", Description: "获取服务运行状态", Group: "base", Method: "GET"},
			// 初始化 API
			&model.Api{Path: "/api/v1/init/postinit", Description: "初始化数据库", Group: "init", Method: "POST"},
			// 用户认证 API
			&model.Api{Path: "/api/v1/auth/postlogin", Description: "用户登录", Group: "auth", Method: "POST"},
		},
	},
	{
		TableName: "sys_users",
		Data: []interface{}{
			&model.User{Username: "admin", Password: util.BcryptHash("123456"), Nickname: "系统管理员", AuthorityId: 1, Authorities: []model.Authority{{AuthorityName: "系统管理员"}}},
		},
	},
}

// CommonDBOperations 定义了数据库操作的公共接口
type CommonDBOperations struct{}

// CreateTable 创建表结构
func (c *CommonDBOperations) CreateTable() error {
	for _, initData := range initialDatas {
		tableName := initData.TableName
		exists := global.DB.Migrator().HasTable(tableName) // 检查表是否存在
		if !exists {
			if err := global.DB.AutoMigrate(initData.Data...); err != nil {
				return fmt.Errorf("创建表 %s 失败: %w", tableName, err)
			}
		}
	}
	return nil
}

func (c *CommonDBOperations) InsertData() error {
	tx := global.DB.Begin() // 回滚事务，避免出现只完成了部分插入的情况。

	for _, initData := range initialDatas {
		for _, data := range initData.Data {
			if initializableData, ok := data.(model.Initializable); ok {
				if err := initializableData.InsertData(global.DB); err != nil {
					tx.Rollback() // 插入失败，回滚事务
					return fmt.Errorf("初始化表 %s 失败: %w", initData.TableName, err)
				}
			} else {
				tx.Rollback() // 类型断言失败，回滚事务
				return fmt.Errorf("数据项 %v 不支持初始化接口", data)
			}

		}
	}
	// 提交事务
	return tx.Commit().Error
}

func ExecuteSQL(dsn string, driver string, sqlStatement string) error {
	db, err := sql.Open(driver, dsn)
	if err != nil {
		return err
	}
	defer db.Close()

	if err = db.Ping(); err != nil {
		return err
	}
	_, err = db.Exec(sqlStatement)
	return err
}
