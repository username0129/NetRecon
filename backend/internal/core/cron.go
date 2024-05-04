package core

import (
	"fmt"
	"github.com/robfig/cron/v3"
	"sync"
	"time"
)

// CronManager 用于管理和存储所有cron任务
type CronManager struct {
	cron      *cron.Cron
	tasks     map[cron.EntryID]string
	tasksLock sync.Mutex
}

// NewTaskManager 创建一个新的任务管理器
func NewTaskManager() *CronManager {
	return &CronManager{
		cron:  cron.New(cron.WithSeconds()), // 使用秒级精度的 cron
		tasks: make(map[cron.EntryID]string),
	}
}

// AddTask 添加新任务并返回任务ID
func (m *CronManager) AddTask(spec string, cmd func()) (cron.EntryID, error) {
	id, err := m.cron.AddFunc(spec, cmd)
	if err != nil {
		return 0, err
	}
	m.tasksLock.Lock()
	m.tasks[id] = spec
	m.tasksLock.Unlock()
	return id, nil
}

// RemoveTask 通过任务ID删除任务
func (m *CronManager) RemoveTask(id cron.EntryID) {
	m.cron.Remove(id)
	m.tasksLock.Lock()
	delete(m.tasks, id)
	m.tasksLock.Unlock()
}

// Start 启动所有cron任务
func (m *CronManager) Start() {
	m.cron.Start()
}

// Stop 停止所有cron任务
func (m *CronManager) Stop() {
	m.cron.Stop()
}

// ListTasks 列出所有任务
func (m *CronManager) ListTasks() map[cron.EntryID]string {
	return m.tasks
}

// RestartTask 重新启动指定ID的任务
func (m *CronManager) RestartTask(id cron.EntryID) (cron.EntryID, error) {
	m.tasksLock.Lock()
	spec, exists := m.tasks[id]
	if !exists {
		m.tasksLock.Unlock()
		return 0, fmt.Errorf("定时任务 ID %v 未找到 ", id)
	}
	m.tasksLock.Unlock()

	// 先删除旧任务
	m.RemoveTask(id)

	// 添加新任务
	newID, err := m.AddTask(spec, func() {
		fmt.Println("任务重新规则: ", time.Now(), ", 运行频率:", spec)
	})
	if err != nil {
		return 0, err
	}
	return newID, nil
}
