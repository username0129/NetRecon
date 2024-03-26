package core

import (
	"backend/internal/global"
	"errors"
	"fmt"
	"net/http"
	"os"
	"time"
)

func StartServer() {
	r := InitializeRout()

	address := fmt.Sprintf("%v:%v", global.Config.System.Ip, global.Config.System.Port)

	global.Logger.Info(fmt.Sprintf("服务端开始监听在: %s", address))

	server := &http.Server{
		Addr:           address,
		Handler:        r,
		ReadTimeout:    20 * time.Second,
		WriteTimeout:   20 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	// 监听并服务
	if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		global.Logger.Error(fmt.Sprintf("服务启动失败: %s", err))
		os.Exit(-1)
	} else {
		global.Logger.Info("服务正常关闭")
		os.Exit(0)
	}
}
