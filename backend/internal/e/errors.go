package e

import "errors"

var (
	ErrDatabaseConfigInvalid  = errors.New("数据库配置文件无效")
	ErrCacheEntryTimeout      = errors.New("缓存条目超时")
	ErrDatabaseNotInitialized = errors.New("数据库未初始化")
	ErrTokenExpired           = errors.New("令牌已过期")
	ErrTokenNotYetActive      = errors.New("令牌尚未激活")
	ErrTokenMalformed         = errors.New("令牌格式错误")
	ErrTokenInvalid           = errors.New("令牌无效")
)
