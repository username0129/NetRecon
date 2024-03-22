package core

import (
	"backend/internal/global"
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

func InitializeViper() *viper.Viper {
	v := viper.New()
	v.SetConfigFile(global.Config.System.ConfigPath) // 设置配置文件路径
	v.SetConfigType("yaml")                          // 设置配置文件类型

	if err := v.ReadInConfig(); err != nil {
		panic(fmt.Errorf("读取配置文件错误：%v\n", err))
	} // 读取配置文件

	v.WatchConfig()
	v.OnConfigChange(func(e fsnotify.Event) {
		global.Logger.Info(fmt.Sprintf("配置文件发生更改：%v", e.Name))
		if err := v.Unmarshal(&global.Config); err != nil {
			panic(fmt.Errorf("配置文件解析错误：%v\n", err))
		}
		// 这里可以添加更多的配置验证逻辑
		global.Logger.Info("配置文件重新加载和解析成功")
	}) // 监视配置文件更改

	if err := v.Unmarshal(&global.Config); err != nil {
		panic(fmt.Errorf("配置文件解析错误：%v\n", err))
	} // 解析配置文件到全局配置

	return v
}
