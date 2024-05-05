package test

import (
	"backend/internal/config"
	"backend/internal/util"
	"fmt"
	"github.com/casbin/casbin/v2/log"
	"testing"
	"time"
)

func TestCron(t *testing.T) {
	manager := config.NewCronManager()
	manager.Start()
	spec := util.TimeToCronSpec(time.Now().Add(1 * time.Minute)) // 3 秒钟之后执行
	fmt.Println(spec)
	// 添加任务
	taskID, err := manager.AddTask(spec, func() {
		fmt.Println("Task executed at", time.Now(), "with spec", spec)
	})
	if err != nil {
		log.LogError(err)
		return
	}
	fmt.Println(taskID)
	time.Sleep(10 * time.Minute)
}
