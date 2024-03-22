package config

// Zap
// @Description: zap 日志配置
type Zap struct {
	Level        string `mapstructure:"level" yaml:"level" json:"level,omitempty"`                            // 最低记录日志等级
	Format       string `mapstructure:"format" yaml:"format" json:"format,omitempty"`                         // 输出格式
	Director     string `mapstructure:"director" yaml:"director" json:"director,omitempty"`                   // 日志保存路径
	EncodeLevel  string `mapstructure:"encode_level" yaml:"encode_level" json:"encode_level,omitempty"`       // 日志级别的编码器
	LogInConsole bool   `mapstructure:"log_in_console" yaml:"log_in_console" json:"log_in_console,omitempty"` // 日志输出到控制台
}
