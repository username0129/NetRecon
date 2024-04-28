package test

import (
	"backend/internal/util"
	"sync"
	"testing"
)

func TestIcmp(t *testing.T) {
	aliveTargets := make(map[string]bool)
	var TargetsMutex sync.Mutex // 互斥锁，保护 Targets
	targets := []string{"192.168.80.2"}
	var wg sync.WaitGroup
	semaphore := make(chan struct{}, 10) // 用于控制并发数量的信号量

	// 首先检测所有目标是否存活
	for _, target := range targets {
		wg.Add(1)
		semaphore <- struct{}{}
		go func(t string) {
			defer wg.Done()
			defer func() { <-semaphore }()
			if alive, _ := util.IcmpCheckAlive(t, 5); alive {
				TargetsMutex.Lock()
				aliveTargets[t] = true
				TargetsMutex.Unlock()
			}
		}(target)
	}
	wg.Wait()
}
