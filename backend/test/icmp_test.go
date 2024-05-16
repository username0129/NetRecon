package test

import (
	"backend/internal/util"
	"fmt"
	"testing"
)

func TestIcmp(t *testing.T) {
	aliveTargets := make(map[string]bool)

	// 解析多个 IP 地址
	addresses, err := util.ParseMultipleIPAddresses("121.37.217.139")
	if err != nil {
		t.Fatalf("解析 IP 地址错误: %v", err)
	}

	// 检测所有目标是否存活
	for _, target := range addresses {
		if alive, err := util.IcmpCheckAlive(target); err != nil {
			fmt.Println(err)
		} else if alive {
			aliveTargets[target] = true
		}
	}

	fmt.Println(aliveTargets)
}
