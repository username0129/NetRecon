package request

import (
	"backend/internal/model"
	"github.com/gofrs/uuid/v5"
)

type CancelTaskRequest struct {
	UUID uuid.UUID `json:"uuid"`
}

type DeleteTaskRequest struct {
	UUID uuid.UUID `json:"uuid"`
}

type DeleteTasksRequest struct {
	UUIDS []uuid.UUID `json:"uuids"`
}

type FetchTasksRequest struct {
	PageInfo
	model.Task
	OrderKey string `json:"orderKey"` // 排序
	Desc     bool   `json:"desc"`     // 排序方式:升序false(默认)|降序true
}

type FetchTaskCountRequest struct {
	TaskType string `json:"taskType"`
}
