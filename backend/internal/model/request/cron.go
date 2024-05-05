package request

type CronAddTaskRequest struct {
	TaskType   string `json:"taskType"`
	Title      string `json:"title"`   // 标题
	Targets    string `json:"targets"` // 目标 IP
	Ports      string `json:"ports"`   // 目标端口
	Timeout    int    `json:"timeout"` // 自定义超时时间
	Threads    int    `json:"threads"` // 线程数
	DictType   string `json:"dictType"`
	CheckAlive bool   `json:"checkAlive"` // 是否进行存活探测
}
