package config

type Mail struct {
	SmtpServer   string `mapstructure:"smtp_server" json:"smtp_server" yaml:"smtp_server"`
	SmtpPort     string `mapstructure:"smtp_port" json:"smtp_port" yaml:"smtp_port"`
	SmtpFrom     string `mapstructure:"smtp_from" json:"smtp_from" yaml:"smtp_from"`
	SmtpPassword string `mapstructure:"smtp_password" json:"smtp_password" yaml:"smtp_password"`
}
