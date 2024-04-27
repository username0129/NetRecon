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
	global.TaskManager[task.UUID] = task.Cancel
	return task, nil
}
