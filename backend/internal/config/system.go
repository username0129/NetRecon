package config

type System struct {
	DBType       string `mapstructure:"db_type" yaml:"db_type" json:"db_type,omitempty"`                   // 数据库类型
	RouterPrefix string `mapstructure:"router_prefix" yaml:"router_prefix" json:"router_prefix,omitempty"` // api 路由前缀
}
