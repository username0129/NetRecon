package controller

import (
	"backend/internal/global"
	"backend/internal/model"
	"backend/internal/model/common"
	"backend/internal/model/request"
	"backend/internal/service"
	"backend/internal/util"
	"errors"
	"fmt"
	"github.com/allegro/bigcache/v3"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
	"strconv"
)

type AuthController struct{}

func (ac *AuthController) PostLogin(c *gin.Context) {
	var logonRequest request.LoginRequest

	if err := c.ShouldBindJSON(&logonRequest); err != nil {
		common.ResponseOk(c, http.StatusBadRequest, "参数解析错误！", nil)
		return
	}

	// 验证码
	openCaptcha := global.Config.Captcha.OpenCaptcha // 是否开启验证码

	key := c.ClientIP() // 客户端 IP

	item, err := global.Cache.Get(key)
	if err != nil {
		if errors.Is(err, bigcache.ErrEntryNotFound) {
			_ = global.Cache.Set(key, []byte("1"))
		} else {
			global.Logger.Error("获取缓存条目错误！", zap.Error(err))
			return
		}
	}
	count, _ := strconv.Atoi(string(item))

	var oc = openCaptcha == 0 || openCaptcha <= count

	if !oc || (logonRequest.CaptchaId != "" && logonRequest.Answer != "" && global.CaptchaStore.Verify(logonRequest.CaptchaId, logonRequest.Answer, true)) {
		u := model.User{Username: logonRequest.Username, Password: logonRequest.Password}
		var user = &model.User{}
		if user, err = service.AuthServiceApp.Login(u); err != nil {
			global.Logger.Error(fmt.Sprintf("用户 %v 登陆失败：%v", logonRequest.Username, err.Error()))
			_ = global.Cache.Set(key, []byte(strconv.Itoa(count+1)))
			common.ResponseOk(c, http.StatusInternalServerError, fmt.Sprintf("登陆失败: %v", err.Error()), nil)
			return
		}
		ac.TokenNext(c, *user) // 用户登录成功，生成 Token
		return
	} else {
		_ = global.Cache.Set(key, []byte(strconv.Itoa(count+1)))
		common.ResponseOk(c, http.StatusUnauthorized, "登陆失败: 验证码错误", nil)
		return
	}
}

func (ac *AuthController) TokenNext(c *gin.Context, user model.User) {
	token, err := util.GenerateJWT(model.CustomClaims{
		UUID:        user.UUID,
		Username:    user.Username,
		AuthorityId: user.AuthorityId,
	})
	if err != nil {
		global.Logger.Error(fmt.Sprintf("用户 %v 登陆失败：获取 token 失败！", user.Username), zap.Error(err))
		common.ResponseOk(c, http.StatusUnauthorized, fmt.Sprintf("用户 %v 登陆失败：获取 token 失败！", user.Username), nil)
	} else {
		global.Logger.Info(fmt.Sprintf("用户 %v 登陆成功！", user.Username))
		common.ResponseOk(c, http.StatusOK, "登陆成功！", gin.H{
			"user":  user,
			"token": token,
		})
	}
}
