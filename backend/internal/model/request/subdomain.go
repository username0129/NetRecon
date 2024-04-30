package request

type SubDomainRequest struct {
	Title    string `json:"title"`    // 标题
	Targets  string `json:"targets"`  // 目标 IP
	Timeout  string `json:"timeout"`  // 自定义超时时间
	Threads  string `json:"threads"`  // 线程数
	DictType string `json:"dictType"` // 字典类型
}
