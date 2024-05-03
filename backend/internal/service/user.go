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

// GetUserInfo  根据 UUID 获取用户详细信息
func (us *UserService) GetUserInfo(userUUID uuid.UUID) (user model.User, err error) {
	if err := global.DB.Model(&model.User{}).Where("uuid = ?", userUUID).First(&user).Error; err != nil {
		return model.User{}, err
	} else {
		return user, nil
	}
}

// GetAdministratorMail 获取管理员邮箱
func (us *UserService) GetAdministratorMail() (mails []string, err error) {
	if err := global.DB.Model(&model.User{}).Select("mail").Where("authority_id = 1").Find(&mails).Error; err != nil {
		return nil, err
	} else {
		return mails, nil
	}
}

// GetUserMailByUUID 根据 UUID 获取用户邮箱
func (us *UserService) GetUserMailByUUID(userUUID uuid.UUID) (mails []string, err error) {
	if err := global.DB.Model(&model.User{}).Select("mail").Where("uuid = ?", userUUID).Find(&mails).Error; err != nil {
		return nil, err
	} else {
		return mails, nil
	}
}
