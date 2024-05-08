package config

import (
	"github.com/robfig/cron/v3"
	"sync"
)

// CronManager 用于管理和存储所有cron任务
type CronManager struct {
	cron      *cron.Cron
	tasks     map[cron.EntryID]string
	tasksLock sync.Mutex
}

// NewCronManager 创建一个新的任务管理器
func NewCronManager() *CronManager {
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
