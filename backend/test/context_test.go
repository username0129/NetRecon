package test

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"
)

func TestContext(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		operation(ctx, 3*time.Second)
	}()

	time.Sleep(2 * time.Second)
	cancel() // 提前取消操作
	wg.Wait()
	//time.Sleep(2 * time.Second) // 给协程足够的时间来响应取消
}

func operation(ctx context.Context, duration time.Duration) {
	select {
	case <-time.After(duration):
		fmt.Println("完成操作")
		return
	case <-ctx.Done():
		fmt.Println("取消操作")
		return
	}
}
