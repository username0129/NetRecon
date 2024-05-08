package request

import (
	"backend/internal/model"
	"github.com/gofrs/uuid/v5"
)

type SubDomainRequest struct {
	Title    string `json:"title"`    // 标题
	Targets  string `json:"targets"`  // 目标 IP
	Timeout  int    `json:"timeout"`  // 自定义超时时间
	Threads  int    `json:"threads"`  // 线程数
	DictType string `json:"dictType"` // 字典类型
}

type FetchSubDomainResultRequest struct {
	PageInfo
	model.SubDomainResult
	OrderKey string `json:"orderKey"` // 排序
	Desc     bool   `json:"desc"`     // 排序方式:升序false(默认)|降序true
}

type DeleteSubDomainResultRequest struct {
	UUID uuid.UUID `json:"uuid"`
}

type DeleteSubDomainResultsRequest struct {
	UUIDS []uuid.UUID `json:"uuids"`
}
