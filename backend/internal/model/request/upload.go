package request

import "backend/internal/model"

type FetchFilesRequest struct {
	PageInfo
	model.File
	OrderKey string `json:"orderKey"` // 排序
	Desc     bool   `json:"desc"`     // 排序方式:升序false(默认)|降序true
}
