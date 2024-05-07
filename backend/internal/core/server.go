package core

import (
	"backend/internal/global"
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"time"
)

func StartServer(ip, port string) {
	for { // 无限循环，以便重启
		r := InitializeRout()

		address := fmt.Sprintf("%v:%v", ip, port)
		global.Logger.Info(fmt.Sprintf("服务端开始监听在: %s", address))

		server := &http.Server{
			Addr:           address,
			Handler:        r,
			ReadTimeout:    20 * time.Second,
			WriteTimeout:   20 * time.Second,
			MaxHeaderBytes: 1 << 20,
		}

		// 启动服务器
		go func() {
			if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
				global.Logger.Error(fmt.Sprintf("服务启动失败: %s", err))
				os.Exit(-1)
			}
		}()

		// 等待重启信号
		<-global.RestartSignal

		// 优雅地关闭服务器
		if err := server.Shutdown(context.Background()); err != nil {
			global.Logger.Error(fmt.Sprintf("服务关闭失败: %s", err))
		}

		global.Logger.Info("服务器正在重启...")
		// 确保有足够的时间来释放资源如端口等
		time.Sleep(5 * time.Second)
	}
}
