package util

import (
	"backend/internal/global"
	"backend/internal/model"
	"context"
	"github.com/gofrs/uuid/v5"
)

// StartNewTask StartNew 初始化新任务并保存到数据库
func StartNewTask(taskType string) (*model.Task, error) {
	ctx, cancel := context.WithCancel(context.Background())
	task := &model.Task{
		UUID:   uuid.Must(uuid.NewV4()),
		Type:   taskType,
		Status: "Running",
		Cancel: cancel,
		Ctx:    ctx,
	}
	if err := task.CreateOrUpdate(); err != nil {
		return nil, err
	}
	global.TaskManager[task.UUID] = task
	return task, nil
}

func CancelTask(uuid uuid.UUID) string {
	if task, exists := global.TaskManager[uuid]; exists {
		task.Cancel() // 调用cancel函数取消任务
		task.Status = "Cancelled"
		return "Task cancelled successfully"
	}
	return "Task not found"
}
