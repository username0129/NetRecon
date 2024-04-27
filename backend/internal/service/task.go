package service

import (
	"backend/internal/global"
	"backend/internal/model"
	"errors"
	"github.com/gofrs/uuid/v5"
)

type TaskService struct{}

var (
	TaskServiceApp = new(TaskService)
)

func (ts TaskService) FetchAllTasks() (tasks []model.Task, err error) {
	if err = global.DB.Find(&tasks).Error; err != nil {
		return nil, err // 返回错误信息
	}
	return tasks, nil // 返回任务列表
}

func (ts TaskService) CancelTask(taskUUID, userUUID uuid.UUID, authorityId uint) (err error) {
	if cancel, exists := global.TaskManager[taskUUID]; exists { // 任务在管理器中存在
		var task model.Task
		if err := global.DB.Model(&task).Where("uuid = ?", taskUUID).First(&task).Error; err != nil { // 查询失败
			return err
		} else {
			if task.Status != "1" { // 任务不正在进行中
				return errors.New("任务不在运行状态")
			}
			if task.CreatorUUID != userUUID || authorityId != 1 { // 任务的发起者不是当前用户 或者 不是管理员
				return errors.New("没有权限取消任务")
			}
			// 取消任务
			cancel()
			if err = task.UpdateStatus("3"); err != nil { // 更新任务状态
				return err
			} else {
				return nil
			}
		}
	} else {
		return errors.New("任务不存在或已经完成")
	}
}
