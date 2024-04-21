package service

import (
	"backend/internal/e"
	"backend/internal/global"
	"backend/internal/model"
	"backend/internal/util"
	"errors"
	"fmt"
	"gorm.io/gorm"
)

type AuthService struct{}

var (
	AuthServiceApp = new(AuthService)
)

func (as *AuthService) Login(u model.User) (userInter *model.User, err error) {
	if global.DB == nil {
		return nil, e.ErrDatabaseNotInitialized
	}

	var user model.User
	if err = global.DB.Model(&model.User{}).Where("username = ?", u.Username).Preload("Authorities").First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("账号或密码错误")
		} else {
			global.Logger.Error(fmt.Sprintf("用户 %v 登陆失败：%v", u.Username, err.Error()))
			return nil, errors.New("出现意料之外的错误，请查看后端日志。")
		}
	} else {
		if ok := util.BcryptCheck(u.Password, user.Password); !ok {
			return nil, errors.New("账号或密码错误")
		} else {
			if user.Enable != 1 {
				return nil, errors.New("账号被冻结")
			}
			return &user, nil
		}
	}
}
