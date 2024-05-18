package test

import (
	"backend/internal/util"
	"fmt"
	"sync"
	"testing"
)

func TestIcmp(t *testing.T) {
	var wg sync.WaitGroup
	semaphore := make(chan struct{}, 100) // 用于控制并发数量的信号量
	addresses, _ := util.ParseMultipleIPAddresses("121.37.217.130-140")
	for _, address := range addresses {
		wg.Add(1)
		go func(t string) {
			defer wg.Done()
			defer func() { <-semaphore }()
			semaphore <- struct{}{}
			alive := util.IcmpCheckAlive(t)
			if alive {
				fmt.Printf("%s is alive\n", t)
			}
		}(address)
	}
	wg.Wait()
}
