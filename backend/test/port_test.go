package test

import (
	"backend/internal/util"
	"backend/lib/gonmap"
	"context"
	"fmt"
	"sync"
	"testing"
	"time"
)

func TestScan(t *testing.T) {
	// 设置目标主机和端口范围
	target := "121.37.217.130"
	list := util.ParsePort("1-65535")

	var wg sync.WaitGroup
	semaphore := make(chan struct{}, 500) // 用于控制并发数量的信号量

	// 创建一个可以取消的context
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// 在2秒后自动取消扫描作为示例
	//time.AfterFunc(1*time.Second, cancel)

	for _, port := range list {
		wg.Add(1)
		semaphore <- struct{}{}
		go func(t string, p int) {
			defer func() {
				<-semaphore
				wg.Done()
			}()
			if err := ctx.Err(); err != nil {
				fmt.Println("新任务添加被阻止，上下文已取消")
				return
			}
			runScan(t, p)
		}(target, port)
	}

	wg.Wait()
	fmt.Println("所有扫描任务已完成或被取消")
}

func runScan(target string, port int) {
	scanner := gonmap.New()
	scanner.SetTimeout(15 * time.Second)
	status, response := scanner.Scan(target, port)
	switch status {
	case gonmap.Closed:
		fmt.Printf("%v:%v %v\n", target, port, "closed")
	// filter 未知状态
	case gonmap.Unknown:
		fmt.Printf("%v:%v %v", target, port, "unknown")
	default:
		fmt.Printf("%v:%v %v ", target, port, "open")
	}
	if response != nil {
		if response.FingerPrint.Service != "" {
			fmt.Printf("Service: %v\n", response.FingerPrint.Service)
		} else {
			fmt.Printf("Service: unknown\n")
		}
	}
}
