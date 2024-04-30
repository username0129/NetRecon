package util

import (
	"os"
	"path/filepath"
)

// GetExecPwd 获取当前执行目录
func GetExecPwd() string {
	exePath, _ := os.Executable()
	exeDir := filepath.Dir(exePath)
	return exeDir
}
