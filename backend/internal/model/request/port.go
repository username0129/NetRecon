package request

type PortScanRequest struct {
	Title      string `json:"title"`      // 标题
	Targets    string `json:"targets"`    // 目标 IP
	Ports      string `json:"ports"`      // 目标端口
	Timeout    string `json:"timeout"`    // 自定义超时时间
	Threads    string `json:"threads"`    // 线程数
	CheckAlive string `json:"checkAlive"` // 是否进行存活探测
}
