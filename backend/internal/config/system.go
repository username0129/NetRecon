package config

type System struct {
	Ip           string `mapstructure:"ip" yaml:"ip" json:"ip,omitempty"`                                  // 后端 IP 地址
	Port         string `mapstructure:"port" yaml:"port" json:"port,omitempty"`                            // 后端端口
	DBType       string `mapstructure:"db_type" yaml:"db_type" json:"db_type,omitempty"`                   // 数据库类型
	ConfigPath   string `mapstructure:"config_path" yaml:"config_path" json:"config_path,omitempty"`       // 配置文件路径
	RouterPrefix string `mapstructure:"router_prefix" yaml:"router_prefix" json:"router_prefix,omitempty"` // api 路由前缀
}
