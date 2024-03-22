package core

import (
	"backend/internal/global"
	"context"
	"github.com/allegro/bigcache/v3"
	"go.uber.org/zap"
	"os"
	"time"
)

// 参考链接：https://liuqh.icu/2021/06/15/go/package/14-bigcache/

func InitializeCache() *bigcache.BigCache {
	cacheConfig := bigcache.DefaultConfig(5 * time.Minute)

	cache, err := bigcache.New(context.Background(), cacheConfig)
	if err != nil {
		global.Logger.Error("初始化 BigCache 失败: %v", zap.Error(err))
		os.Exit(-1)
	}

	return cache
}
