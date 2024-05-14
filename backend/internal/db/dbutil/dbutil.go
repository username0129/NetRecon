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
			// 测试模块
			&model.CasbinRule{Ptype: "p", V0: "1", V1: "/api/v1/base/gethealth", V2: "GET"}, // 检查路由
			&model.CasbinRule{Ptype: "p", V0: "2", V1: "/api/v1/base/gethealth", V2: "GET"}, //

			// 管理员模块
			&model.CasbinRule{Ptype: "p", V0: "1", V1: "/api/v1/user/postdeleteuserinfo", V2: "POST"},     // 删除用户信息
			&model.CasbinRule{Ptype: "p", V0: "1", V1: "/api/v1/user/postadduserinfo", V2: "POST"},        // 添加用户信息
			&model.CasbinRule{Ptype: "p", V0: "1", V1: "/api/v1/user/postfetchusers", V2: "POST"},         // 获取所有用户信息
			&model.CasbinRule{Ptype: "p", V0: "1", V1: "/api/v1/operation/postfetchresult", V2: "POST"},   // 获取操作记录
			&model.CasbinRule{Ptype: "p", V0: "1", V1: "/api/v1/operation/postdeleteresults", V2: "POST"}, // 获取操作记录

			// 路由模块
			&model.CasbinRule{Ptype: "p", V0: "1", V1: "/api/v1/route/getroute", V2: "GET"}, // 根据当前用户角色获取路由信息
			&model.CasbinRule{Ptype: "p", V0: "2", V1: "/api/v1/route/getroute", V2: "GET"},

			// 用户信息模块
			&model.CasbinRule{Ptype: "p", V0: "1", V1: "/api/v1/user/getuserinfo", V2: "GET"}, // 根据 Token 获取当前用户信息
			&model.CasbinRule{Ptype: "p", V0: "2", V1: "/api/v1/user/getuserinfo", V2: "GET"},
			&model.CasbinRule{Ptype: "p", V0: "1", V1: "/api/v1/user/postupdateuserinfo", V2: "POST"}, // 更新个人信息
			&model.CasbinRule{Ptype: "p", V0: "2", V1: "/api/v1/user/postupdateuserinfo", V2: "POST"},
			&model.CasbinRule{Ptype: "p", V0: "1", V1: "/api/v1/user/postresetpassword", V2: "POST"}, // 重置密码
			&model.CasbinRule{Ptype: "p", V0: "2", V1: "/api/v1/user/postresetpassword", V2: "POST"},
			&model.CasbinRule{Ptype: "p", V0: "1", V1: "/api/v1/user/postupdatepassword", V2: "POST"}, // 更新用户密码
			&model.CasbinRule{Ptype: "p", V0: "2", V1: "/api/v1/user/postupdatepassword", V2: "POST"},

			// 端口扫描模块
			&model.CasbinRule{Ptype: "p", V0: "1", V1: "/api/v1/portscan/postportscan", V2: "POST"}, // 提交端口扫描任务
			&model.CasbinRule{Ptype: "p", V0: "2", V1: "/api/v1/portscan/postportscan", V2: "POST"},
			&model.CasbinRule{Ptype: "p", V0: "1", V1: "/api/v1/portscan/postfetchresult", V2: "POST"}, // 查询结果
			&model.CasbinRule{Ptype: "p", V0: "2", V1: "/api/v1/portscan/postfetchresult", V2: "POST"},
			&model.CasbinRule{Ptype: "p", V0: "1", V1: "/api/v1/portscan/postdeleteresult", V2: "POST"}, // 删除扫描结果
			&model.CasbinRule{Ptype: "p", V0: "2", V1: "/api/v1/portscan/postdeleteresult", V2: "POST"},
			&model.CasbinRule{Ptype: "p", V0: "1", V1: "/api/v1/portscan/postdeleteresults", V2: "POST"}, // 批量删除扫描结果
			&model.CasbinRule{Ptype: "p", V0: "2", V1: "/api/v1/portscan/postdeleteresults", V2: "POST"},
			&model.CasbinRule{Ptype: "p", V0: "1", V1: "/api/v1/portscan/postaddresult", V2: "POST"}, // 咱家扫描结果
			&model.CasbinRule{Ptype: "p", V0: "2", V1: "/api/v1/portscan/postaddresult", V2: "POST"},
			&model.CasbinRule{Ptype: "p", V0: "1", V1: "/api/v1/portscan/postexportdata", V2: "POST"}, // 导出数据
			&model.CasbinRule{Ptype: "p", V0: "2", V1: "/api/v1/portscan/postexportdata", V2: "POST"},

			// 任务管理模块
			&model.CasbinRule{Ptype: "p", V0: "1", V1: "/api/v1/task/postcanceltask", V2: "POST"}, // 根据任务 UUID 取消任务
			&model.CasbinRule{Ptype: "p", V0: "2", V1: "/api/v1/task/postcanceltask", V2: "POST"},
			&model.CasbinRule{Ptype: "p", V0: "1", V1: "/api/v1/task/postfetchtasks", V2: "POST"}, // 根据任务 UUID 获取数据
			&model.CasbinRule{Ptype: "p", V0: "2", V1: "/api/v1/task/postfetchtasks", V2: "POST"},
			&model.CasbinRule{Ptype: "p", V0: "1", V1: "/api/v1/task/postdeletetask", V2: "POST"}, // 根据任务 UUID 删除任务
			&model.CasbinRule{Ptype: "p", V0: "2", V1: "/api/v1/task/postdeletetask", V2: "POST"},
			&model.CasbinRule{Ptype: "p", V0: "1", V1: "/api/v1/task/postdeletetasks", V2: "POST"}, // 批量根据任务 UUID 删除任务
			&model.CasbinRule{Ptype: "p", V0: "2", V1: "/api/v1/task/postdeletetasks", V2: "POST"},
			&model.CasbinRule{Ptype: "p", V0: "1", V1: "/api/v1/task/postfetchtaskcount", V2: "POST"},
			&model.CasbinRule{Ptype: "p", V0: "2", V1: "/api/v1/task/postfetchtaskcount", V2: "POST"},

			// 子域名爆破模块
			&model.CasbinRule{Ptype: "p", V0: "1", V1: "/api/v1/subdomain/postbrutesubdomains", V2: "POST"}, // 提交子域名扫描任务
			&model.CasbinRule{Ptype: "p", V0: "2", V1: "/api/v1/subdomain/postbrutesubdomains", V2: "POST"},
			&model.CasbinRule{Ptype: "p", V0: "1", V1: "/api/v1/subdomain/postfetchresult", V2: "POST"}, // 查询结果
			&model.CasbinRule{Ptype: "p", V0: "2", V1: "/api/v1/subdomain/postfetchresult", V2: "POST"},
			&model.CasbinRule{Ptype: "p", V0: "1", V1: "/api/v1/subdomain/postdeleteresult", V2: "POST"}, // 删除扫描结果
			&model.CasbinRule{Ptype: "p", V0: "2", V1: "/api/v1/subdomain/postdeleteresult", V2: "POST"},
			&model.CasbinRule{Ptype: "p", V0: "1", V1: "/api/v1/subdomain/postdeleteresults", V2: "POST"}, // 批量删除扫描结果
			&model.CasbinRule{Ptype: "p", V0: "2", V1: "/api/v1/subdomain/postdeleteresults", V2: "POST"},
			&model.CasbinRule{Ptype: "p", V0: "1", V1: "/api/v1/subdomain/postaddresult", V2: "POST"}, // 咱家扫描结果
			&model.CasbinRule{Ptype: "p", V0: "2", V1: "/api/v1/subdomain/postaddresult", V2: "POST"},
			&model.CasbinRule{Ptype: "p", V0: "1", V1: "/api/v1/subdomain/postexportdata", V2: "POST"}, // 导出数据
			&model.CasbinRule{Ptype: "p", V0: "2", V1: "/api/v1/subdomain/postexportdata", V2: "POST"},

			// 资产管理模块
			&model.CasbinRule{Ptype: "p", V0: "1", V1: "/api/v1/asset/postaddasset", V2: "POST"}, // 添加资产
			&model.CasbinRule{Ptype: "p", V0: "2", V1: "/api/v1/asset/postaddasset", V2: "POST"},
			&model.CasbinRule{Ptype: "p", V0: "1", V1: "/api/v1/asset/postdeleteasset", V2: "POST"}, // 删除资产
			&model.CasbinRule{Ptype: "p", V0: "2", V1: "/api/v1/asset/postdeleteasset", V2: "POST"},
			&model.CasbinRule{Ptype: "p", V0: "1", V1: "/api/v1/asset/postdeleteassets", V2: "POST"}, // 删除资产
			&model.CasbinRule{Ptype: "p", V0: "2", V1: "/api/v1/asset/postdeleteassets", V2: "POST"},
			&model.CasbinRule{Ptype: "p", V0: "1", V1: "/api/v1/asset/postupdateasset", V2: "POST"}, // 修改资产
			&model.CasbinRule{Ptype: "p", V0: "2", V1: "/api/v1/asset/postupdateasset", V2: "POST"},
			&model.CasbinRule{Ptype: "p", V0: "1", V1: "/api/v1/asset/postfetchasset", V2: "POST"}, // 获取
			&model.CasbinRule{Ptype: "p", V0: "2", V1: "/api/v1/asset/postfetchasset", V2: "POST"},
			&model.CasbinRule{Ptype: "p", V0: "2", V1: "/api/v1/asset/postfetchassetipdetail", V2: "POST"}, // 获取资产 IP 详情
			&model.CasbinRule{Ptype: "p", V0: "2", V1: "/api/v1/asset/postfetchassetipdetail", V2: "POST"},
			&model.CasbinRule{Ptype: "p", V0: "2", V1: "/api/v1/asset/postfetchassetdomaindetail", V2: "POST"}, // 获取资产域名详情
			&model.CasbinRule{Ptype: "p", V0: "2", V1: "/api/v1/asset/postfetchassetdomaindetail", V2: "POST"},

			// 计划任务模块
			&model.CasbinRule{Ptype: "p", V0: "1", V1: "/api/v1/cron/postaddtask", V2: "POST"}, // 添加计划任务
			&model.CasbinRule{Ptype: "p", V0: "2", V1: "/api/v1/cron/postaddtask", V2: "POST"},

			// fofa 任务管理
			&model.CasbinRule{Ptype: "p", V0: "1", V1: "/api/v1/fofa/postfofasearch", V2: "POST"}, // 添加 Fofa 查询
			&model.CasbinRule{Ptype: "p", V0: "2", V1: "/api/v1/fofa/postfofasearch", V2: "POST"},
			&model.CasbinRule{Ptype: "p", V0: "1", V1: "/api/v1/fofa/postexportdata", V2: "POST"},
			&model.CasbinRule{Ptype: "p", V0: "2", V1: "/api/v1/fofa/postexportdata", V2: "POST"},

			// 图片管理
			&model.CasbinRule{Ptype: "p", V0: "1", V1: "/api/v1/upload/postuploadfile", V2: "POST"}, // 文件上传
			&model.CasbinRule{Ptype: "p", V0: "2", V1: "/api/v1/upload/postuploadfile", V2: "POST"},
			&model.CasbinRule{Ptype: "p", V0: "1", V1: "/api/v1/upload/postdeletefile", V2: "POST"}, // 删除文件
			&model.CasbinRule{Ptype: "p", V0: "2", V1: "/api/v1/upload/postdeletefile", V2: "POST"},
			&model.CasbinRule{Ptype: "p", V0: "1", V1: "/api/v1/upload/postfetchfiles", V2: "POST"}, // 获取文件列表
			&model.CasbinRule{Ptype: "p", V0: "2", V1: "/api/v1/upload/postfetchfiles", V2: "POST"},
		},
	},
	{
		TableName: "sys_routes",
		Data: []interface{}{
			// 顶级菜单
			&model.Route{ParentId: 0, Meta: model.Meta{Title: "仪表盘", Icon: "odometer"}, Name: "Dashboard", Path: "dashboard", Component: "views/dashboard/IndexView.vue", Authorities: []model.Authority{{AuthorityName: "系统管理员"}, {AuthorityName: "普通用戶"}}},
			&model.Route{ParentId: 0, Meta: model.Meta{Title: "管理面板", Icon: "MessageBox"}, Name: "Admin", Path: "admin", Component: "views/admin/IndexView.vue", Authorities: []model.Authority{{AuthorityName: "系统管理员"}}},
			&model.Route{ParentId: 0, Meta: model.Meta{Title: "子域名扫描", Icon: "Coordinate"}, Name: "Subdomain", Path: "subdomain", Component: "views/subdomain/IndexView.vue", Authorities: []model.Authority{{AuthorityName: "系统管理员"}, {AuthorityName: "普通用戶"}}},
			&model.Route{ParentId: 0, Meta: model.Meta{Hidden: true, Title: "子域名扫描详情", Icon: "avatar"}, Name: "SubdomainDetail", Path: "subdomaindetail", Component: "views/subdomain/details/IndexView.vue", Authorities: []model.Authority{{AuthorityName: "系统管理员"}, {AuthorityName: "普通用戶"}}},
			&model.Route{ParentId: 0, Meta: model.Meta{Title: "IP 端口扫描", Icon: "Cloudy"}, Name: "PortScan", Path: "portscan", Component: "views/portscan/IndexView.vue", Authorities: []model.Authority{{AuthorityName: "系统管理员"}, {AuthorityName: "普通用戶"}}},
			&model.Route{ParentId: 0, Meta: model.Meta{Hidden: true, Title: "IP 端口扫描详情", Icon: "avatar"}, Name: "PortscanDetail", Path: "portscandetail", Component: "views/portscan/details/IndexView.vue", Authorities: []model.Authority{{AuthorityName: "系统管理员"}, {AuthorityName: "普通用戶"}}},
			&model.Route{ParentId: 0, Meta: model.Meta{Title: "资产管理", Icon: "Reading"}, Name: "Asset", Path: "asset", Component: "views/asset/IndexView.vue", Authorities: []model.Authority{{AuthorityName: "系统管理员"}, {AuthorityName: "普通用戶"}}},
			&model.Route{ParentId: 0, Meta: model.Meta{Title: "计划任务", Icon: "Coin"}, Name: "Cron", Path: "cron", Component: "views/cron/IndexView.vue", Authorities: []model.Authority{{AuthorityName: "系统管理员"}, {AuthorityName: "普通用戶"}}},
			&model.Route{ParentId: 0, Meta: model.Meta{Hidden: true, Title: "个人信息", Icon: "message"}, Name: "Person", Path: "person", Component: "views/person/IndexView.vue", Authorities: []model.Authority{{AuthorityName: "系统管理员"}, {AuthorityName: "普通用戶"}}},
			&model.Route{ParentId: 0, Meta: model.Meta{Title: "FOFA 任务", Icon: "Rank"}, Name: "Fofa", Path: "Fofa", Component: "views/fofa/IndexView.vue", Authorities: []model.Authority{{AuthorityName: "系统管理员"}, {AuthorityName: "普通用戶"}}},
			&model.Route{ParentId: 0, Meta: model.Meta{Title: "图片管理", Icon: "Notebook"}, Name: "Upload", Path: "upload", Component: "views/upload/IndexView.vue", Authorities: []model.Authority{{AuthorityName: "系统管理员"}}},

			// 管理员菜单
			&model.Route{ParentId: 2, Meta: model.Meta{Title: "用户管理", Icon: "coordinate"}, Name: "User", Path: "user", Component: "views/admin/user/IndexView.vue", Authorities: []model.Authority{{AuthorityName: "系统管理员"}}},
			&model.Route{ParentId: 2, Meta: model.Meta{Title: "操作历史", Icon: "pie-chart"}, Name: "Operation", Path: "operation", Component: "views/admin/operation/IndexView.vue", Authorities: []model.Authority{{AuthorityName: "系统管理员"}}},

			// 任务管理
			//&model.Route{ParentId: 2, Meta: model.Meta{Title: "从企业名收集资产", Icon: "avatar"}, Name: "Authority", Path: "authority", Component: "views/admin/AuthorityView.vue", Authorities: []model.Authority{{AuthorityName: "系统管理员"}}},
		},
	},
	{
		TableName: "sys_users",
		Data: []interface{}{
			&model.User{Username: "admin", Password: util.BcryptHash("123456"), Nickname: "系统管理员", AuthorityId: "1", Mail: "q.nah@outlook.com"},
			&model.User{Username: "guest", Password: util.BcryptHash("guest"), Nickname: "测试账户", AuthorityId: "2", Mail: "q.nah@outlook.com"},
		},
	},
	{
		TableName: "sys_tasks",
		Data: []interface{}{
			&model.Task{},
		},
	},
	{
		TableName: "sys_port_scan_result",
		Data: []interface{}{
			&model.PortScanResult{},
		},
	},
	{
		TableName: "sys_subdomain_results",
		Data: []interface{}{
			&model.SubDomainResult{},
		},
	},
	{
		TableName: "sys_network_assets",
		Data: []interface{}{
			&model.Asset{},
		},
	},
	{
		TableName: "sys_operation_records",
		Data: []interface{}{
			&model.OperationRecord{},
		},
	},
	{
		TableName: "sys_files",
		Data: []interface{}{
			&model.File{},
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
