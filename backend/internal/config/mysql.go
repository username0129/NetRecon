package config

import (
	"backend/internal/e"
	"fmt"
	"time"
)

// Mysql
// @Description: MySQL 配置文件
type Mysql struct {
	Host            string        `mapstructure:"host" yaml:"host" json:"host,omitempty"`                                        // 数据库主机地址
	Port            string        `mapstructure:"port" yaml:"port" json:"port,omitempty"`                                        // 数据库所用端口
	DBName          string        `mapstructure:"db_name" yaml:"db_name" json:"db_name,omitempty"`                               // 数据库名
	Username        string        `mapstructure:"username" yaml:"username" json:"username,omitempty"`                            // 数据库用户名
	Password        string        `mapstructure:"password" yaml:"password" json:"password,omitempty"`                            // 数据库密码
	MaxIdleConns    int           `mapstructure:"max_idle_conns" yaml:"max_idle_conns" json:"max_idle_conns,omitempty"`          // 空闲连接最大数量
	MaxOpenConns    int           `mapstructure:"max_open_conns" yaml:"max_open_conns" json:"max_open_conns,omitempty"`          // 打开连接的最大数量
	ConnMaxLifetime time.Duration `mapstructure:"conn_max_lifetime" yaml:"conn_max_lifetime" json:"conn_max_lifetime,omitempty"` // 空闲连接存活时间
	Config          string        `mapstructure:"config" yaml:"config" json:"config"`                                            // 连接配置
	LogMode         string        `mapstructure:"log_mode" yaml:"log_mode" json:"log_mode"`                                      // 是否开启 gorm 全局日志
	LogZap          bool          `mapstructure:"log_zap" yaml:"log_zap" json:"log_zap"`                                         // 是否打印日志到 zap
}

func (m *Mysql) Check() error {
	if m.Host == "" {
		m.Host = "127.0.0.1"
	}
	if m.Port == "" {
		m.Port = "3306"
	}
	if m.Username == "" || m.DBName == "" {
		return e.ErrDatabaseConfigInvalid
	}
	return nil
}

// GetDSN
//
//	@Description: 获取用于连接 MySQL 数据库的 DSN
//	@receiver mysqlConfig
//	@return string
func (m *Mysql) GetDSN() string {
	return fmt.Sprintf("%v:%v@tcp(%v:%v)/%v?%v", m.Username, m.Password, m.Host, m.Port, m.DBName, m.Config)
}

// GetEmptyDSN
//
//	@Description: 获取用于连接 MySQL 的 DSN
//	@receiver mysqlConfig
//	@return string
func (m *Mysql) GetEmptyDSN() string {
	return fmt.Sprintf("%v:%v@tcp(%v:%v)/", m.Username, m.Password, m.Host, m.Port)
}
