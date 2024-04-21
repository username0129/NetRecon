package service

import (
	"backend/internal/global"
	"backend/internal/model"
	"github.com/gofrs/uuid/v5"
)

type UserService struct{}

var (
	UserServiceApp = new(UserService)
)

// GetUserInfo
// 根据 UUID 获取用户详细信息
func (us *UserService) GetUserInfo(uuid uuid.UUID) (user model.User, err error) {
	if err := global.DB.Model(&model.User{}).Where("uuid = ?", uuid).Preload("Authorities").First(&user).Error; err != nil {
		return model.User{}, err
	} else {
		return user, nil
	}
}
