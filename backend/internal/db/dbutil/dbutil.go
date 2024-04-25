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
			&model.Authority{AuthorityName: "普通用戶"},
		},
	},
	{
		TableName: "casbin_role",
		Data: []interface{}{
			// 管理员组
			&model.CasbinRule{Ptype: "p", V0: "1", V1: "/api/v1/user/getuserinfo", V2: "GET"},
			&model.CasbinRule{Ptype: "p", V0: "1", V1: "/api/v1/user/postuserinfo", V2: "POST"},
			&model.CasbinRule{Ptype: "p", V0: "1", V1: "/api/v1/route/getroute", V2: "GET"},
			&model.CasbinRule{Ptype: "p", V0: "1", V1: "/api/v1/portscan/postportscan", V2: "POST"},

			// 普通用户组
			&model.CasbinRule{Ptype: "p", V0: "1", V1: "/api/v1/user/getuserinfo", V2: "GET"},
			&model.CasbinRule{Ptype: "p", V0: "2", V1: "/api/v1/user/postuserinfo", V2: "POST"},
			&model.CasbinRule{Ptype: "p", V0: "2", V1: "/api/v1/route/getroute", V2: "GET"},
			&model.CasbinRule{Ptype: "p", V0: "2", V1: "/api/v1/portscan/postportscan", V2: "POST"},
		},
	},
	{
		TableName: "sys_routes",
		Data: []interface{}{
			// 顶级菜单
			&model.Route{ParentId: 0, Meta: model.Meta{Title: "仪表盘", Icon: "odometer"}, Name: "Dashboard", Path: "dashboard", Component: "views/dashboard/IndexView.vue", Authorities: []model.Authority{{AuthorityName: "系统管理员"}, {AuthorityName: "普通用戶"}}},
			&model.Route{ParentId: 0, Meta: model.Meta{Title: "管理面板", Icon: "user"}, Name: "Admin", Path: "admin", Component: "views/admin/IndexView.vue", Authorities: []model.Authority{{AuthorityName: "系统管理员"}}},
			&model.Route{ParentId: 0, Meta: model.Meta{Title: "任务管理", Icon: "paperclip"}, Name: "Task", Path: "task", Component: "views/task/TaskView.vue", Authorities: []model.Authority{{AuthorityName: "系统管理员"}, {AuthorityName: "普通用戶"}}},
			&model.Route{ParentId: 0, Meta: model.Meta{Title: "个人信息", Icon: "message"}, Name: "Person", Path: "person", Component: "views/person/IndexView.vue", Authorities: []model.Authority{{AuthorityName: "系统管理员"}, {AuthorityName: "普通用戶"}}},

			// 管理员菜单
			&model.Route{ParentId: 2, Meta: model.Meta{Title: "角色管理", Icon: "avatar"}, Name: "Authority", Path: "authority", Component: "views/admin/AuthorityView.vue", Authorities: []model.Authority{{AuthorityName: "系统管理员"}}},
			&model.Route{ParentId: 2, Meta: model.Meta{Title: "用户管理", Icon: "coordinate"}, Name: "User", Path: "user", Component: "views/admin/UserView.vue", Authorities: []model.Authority{{AuthorityName: "系统管理员"}}},
			&model.Route{ParentId: 2, Meta: model.Meta{Title: "操作历史", Icon: "pie-chart"}, Name: "Operation", Path: "operation", Component: "views/admin/OperationView.vue", Authorities: []model.Authority{{AuthorityName: "系统管理员"}}},

			// 任务管理
			&model.Route{ParentId: 2, Meta: model.Meta{Title: "子域名收集", Icon: "avatar"}, Name: "Authority", Path: "authority", Component: "views/admin/AuthorityView.vue", Authorities: []model.Authority{{AuthorityName: "系统管理员"}}},
			&model.Route{ParentId: 2, Meta: model.Meta{Title: "IP 端口扫描", Icon: "avatar"}, Name: "Authority", Path: "authority", Component: "views/admin/AuthorityView.vue", Authorities: []model.Authority{{AuthorityName: "系统管理员"}}},
			&model.Route{ParentId: 2, Meta: model.Meta{Title: "FOFA 任务下发", Icon: "avatar"}, Name: "Authority", Path: "authority", Component: "views/admin/AuthorityView.vue", Authorities: []model.Authority{{AuthorityName: "系统管理员"}}},
			&model.Route{ParentId: 2, Meta: model.Meta{Title: "从企业名收集资产", Icon: "avatar"}, Name: "Authority", Path: "authority", Component: "views/admin/AuthorityView.vue", Authorities: []model.Authority{{AuthorityName: "系统管理员"}}},
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
			&model.User{Username: "admin", Password: util.BcryptHash("123456"), Nickname: "系统管理员", AuthorityId: 1},
			&model.User{Username: "guest", Password: util.BcryptHash("guest"), Nickname: "测试账户", AuthorityId: 2},
		},
	},
	{
		TableName: "sys_port_scan_result",
		Data: []interface{}{
			&model.PortScanResult{Open: false},
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
