package request

import (
	"backend/internal/model"
	"github.com/gofrs/uuid/v5"
)

type AddAssetRequest struct {
	Title   string `json:"title"`
	Domains string `json:"domains"`
	IPs     string `json:"ips"`
}

type DeleteAssetRequest struct {
	UUID uuid.UUID `json:"uuid"`
}

type UpdateAssetRequest struct {
	UUID    uuid.UUID `json:"uuid"`
	Title   string    `json:"title"`
	Domains string    `json:"domains"`
	IPs     string    `json:"ips"`
}

type FetchAssetsRequest struct {
	PageInfo
	model.Asset
	OrderKey string `json:"orderKey"` // 排序
	Desc     bool   `json:"desc"`     // 排序方式:升序false(默认)|降序true
}
