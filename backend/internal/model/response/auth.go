package response

import "backend/internal/model"

type LoginResponse struct {
	User  model.User `json:"user"`  // 用户信息
	Token string     `json:"token"` // 登陆密码
}
