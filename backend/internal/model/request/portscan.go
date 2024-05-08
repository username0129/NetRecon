package request

import (
	"backend/internal/model"
	"github.com/gofrs/uuid/v5"
)

type FetchPortScanResultRequest struct {
	PageInfo
	model.PortScanResult
	OrderKey string `json:"orderKey"` // 排序
	Desc     bool   `json:"desc"`     // 排序方式:升序false(默认)|降序true
}

type DeletePortScanResultRequest struct {
	UUID uuid.UUID `json:"uuid"`
}

type DeletePortScanResultsRequest struct {
	UUIDS []uuid.UUID `json:"uuids"`
}
type PortScanRequest struct {
	Title      string `json:"title"`   // 标题
	Targets    string `json:"targets"` // 目标 IP
	Ports      string `json:"ports"`   // 目标端口
	Timeout    int    `json:"timeout"` // 自定义超时时间
	Threads    int    `json:"threads"` // 线程数
	DictType   string `json:"dictType"`
	CheckAlive bool   `json:"checkAlive"` // 是否进行存活探测
}
