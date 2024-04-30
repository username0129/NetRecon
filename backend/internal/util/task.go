package util

import (
	"backend/internal/global"
	"backend/internal/model"
	"context"
	"github.com/gofrs/uuid/v5"
)

// StartNewTask StartNew 初始化新任务并保存到数据库
func StartNewTask(title, taskType, targets string, userUUID uuid.UUID) (*model.Task, error) {
	ctx, cancel := context.WithCancel(context.Background())
	task := &model.Task{
		Title:       title,
		UUID:        uuid.Must(uuid.NewV4()),
		Targets:     targets,
		CreatorUUID: userUUID,
		Type:        taskType,
		Status:      "1",
		Cancel:      cancel,
		Ctx:         ctx,
	}
	if err := task.CreateOrUpdate(); err != nil {
		return nil, err
	}
	global.TaskManager[task.UUID] = task.Cancel
	return task, nil
}
