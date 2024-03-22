package util

import (
	"os"
)

// IsPathExist IsPathExists
//
//	@Description: 判断文件夹是否存在
//	@param path
//	@return bool
//	@return error
func IsPathExist(path string) (bool, error) {
	if fi, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			return false, nil
		} // 路径不存在
		return false, err // 其他错误，如权限问题
	} else {
		return fi.IsDir(), nil // 目录返回 true
	}
}
