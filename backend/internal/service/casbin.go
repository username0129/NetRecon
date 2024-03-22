package service

import (
	"backend/internal/global"
	"github.com/casbin/casbin/v2"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"sync"
	"time"
)

type CasbinService struct{}

var CasbinServiceApp = new(CasbinService)

var (
	once     sync.Once
	enforcer *casbin.SyncedCachedEnforcer
)

// GetCasbin
//
//	@Description: 从数据库读取 Casbin 配置信息
//	@receiver cs
func (cs *CasbinService) GetCasbin() *casbin.SyncedCachedEnforcer {
	once.Do(func() {
		adapter, err := gormadapter.NewAdapterByDB(global.DB)
		if err != nil {
			return
		}
		enforcer, err = casbin.NewSyncedCachedEnforcer("./config/rbac_model.conf", adapter) // 带有缓存和同步，提高性能以及保证数据同步
		if err != nil {
			return
		}
		enforcer.SetExpireTime(3600 * time.Second)
		_ = enforcer.LoadPolicy()
	})
	return enforcer
}
