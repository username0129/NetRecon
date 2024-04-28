package test

import (
	"backend/internal/util"
	"fmt"
	"testing"
)

func TestIcmp(t *testing.T) {
	// 示例调用
	alive, err := util.IcmpCheckAlive("192.168.1.3", 5)
	if err != nil {
		fmt.Println("Error:", err)
	} else if alive {
		fmt.Println("目标存活")
	} else {
		fmt.Println("目标不存活")
	}
}
