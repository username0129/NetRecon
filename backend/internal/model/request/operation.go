package request

import (
	"backend/internal/model"
	"github.com/gofrs/uuid/v5"
)

type FetchOperationRequest struct {
	PageInfo
	model.OperationRecord
	Username string `json:"username"`
	OrderKey string `json:"orderKey"` // 排序
	Desc     bool   `json:"desc"`     // 排序方式:升序false(默认)|降序true
}

type DeleteOperationsRequest struct {
	UUIDS []uuid.UUID `json:"uuids"`
}
