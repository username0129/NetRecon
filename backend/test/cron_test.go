package test

import (
	"fmt"
	"github.com/robfig/cron/v3"
	"testing"
	"time"
)

func TestCron(t *testing.T) {
	c := cron.New(cron.WithSeconds()) // 支持到秒的精度

	// 每天上午10:15执行
	c.AddFunc("0 15 10 * * *", func() {
		fmt.Println("Task is executed.", time.Now())
	})

	c.Start()

	// 让程序持续运行
	select {}
}
