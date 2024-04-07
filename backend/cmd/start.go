package cmd

import (
	"backend/internal/core"
	"backend/internal/global"
	"backend/internal/logger"
	"github.com/spf13/cobra"
)

var (
	configPath string // 配置文件路径
	ip         string // 后端监听 IP
	port       string // 后端监听端口
	startCmd   = &cobra.Command{
		Use:   "start",
		Short: "Start the Gin web server",
		Run: func(cmd *cobra.Command, args []string) {
			start()
		},
	}
)

func init() {
	startCmd.Flags().StringVarP(&configPath, "config", "c", "./config/config.yaml", "配置文件路径")
	startCmd.Flags().StringVarP(&ip, "ip", "i", "0.0.0.0", "后端 IP 地址")
	startCmd.Flags().StringVarP(&port, "port", "p", "8080", "后端监听地址")
}

func start() {
	global.Viper = core.InitializeViper(configPath) // 初始化并加载 Viper
	global.Logger = logger.InitializeLogger()       // 初始化 Zap 日志
	global.Cache = core.InitializeCache()           // 初始化 BigCache
	global.DB = core.InitializeDB()                 // 获取数据库连接
	if global.DB != nil {
		db, _ := global.DB.DB()
		defer db.Close()
	}
	core.StartServer(ip, port)
}
