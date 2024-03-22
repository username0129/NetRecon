package config

// Server
// @Description: 后端配置
type Server struct {
	// 系统配置
	Jwt     Jwt     `mapstructure:"jwt" yaml:"jwt" json:"jwt"`             // jwt 认证配置
	System  System  `mapstructure:"system" yaml:"system" json:"system"`    // 系统配置
	Zap     Zap     `mapstructure:"zap" yaml:"zap" json:"zap"`             // zap 日志配置
	Captcha Captcha `mapstructure:"captcha" json:"captcha" yaml:"captcha"` // 验证码配置

	// 数据库配置
	Mysql Mysql `mapstructure:"mysql" yaml:"mysql" json:"mysql"` // mysql 配置
}
