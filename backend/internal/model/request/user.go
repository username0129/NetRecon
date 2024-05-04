package request

import (
	"backend/internal/model"
	"github.com/gofrs/uuid/v5"
)

type FetchUsersRequest struct {
	PageInfo
	model.User
	OrderKey string `json:"orderKey"` // 排序
	Desc     bool   `json:"desc"`     // 排序方式:升序false(默认)|降序true
}

type UpdateUserRequest struct {
	UUID        uuid.UUID `json:"uuid"`
	Username    string    `json:"username"`
	Nickname    string    `json:"nickname"`
	AuthorityId string    `json:"authorityId"`
	Mail        string    `json:"mail"`
	Enable      string    `json:"enable"`
}

type AddUserRequest struct {
	Username    string `json:"username"`
	Password    string `json:"password"`
	Nickname    string `json:"nickname"`
	AuthorityId string `json:"authorityId"`
	Mail        string `json:"mail"`
	Enable      string `json:"enable"`
}

type DeleteUserRequest struct {
	UUID uuid.UUID `json:"uuid"`
}

type ResetPasswordRequest struct {
	UUID uuid.UUID `json:"uuid"`
}
