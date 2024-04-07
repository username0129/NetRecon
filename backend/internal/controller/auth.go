package controller

import (
	"backend/internal/global"
	"backend/internal/model"
	"backend/internal/model/response"
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
	var logonRequest model.LoginRequest

	if err := c.ShouldBindJSON(&logonRequest); err != nil {
		response.Response(c, http.StatusInternalServerError, "参数解析错误！", nil)
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
			response.Response(c, http.StatusInternalServerError, fmt.Sprintf("用户 %v 登陆失败：%v", logonRequest.Username, err.Error()), nil)
			return
		} // 用户身份校验失败
		if user.Enable != 1 {
			global.Logger.Error(fmt.Sprintf("用户 %v 登陆失败：用户被冻结，禁止登录!", logonRequest.Username))
			_ = global.Cache.Set(key, []byte(strconv.Itoa(count+1)))
			response.Response(c, http.StatusInternalServerError, fmt.Sprintf("用户 %v 登陆失败：用户被冻结，禁止登录！", logonRequest.Username), nil)
			return
		} // 用户账户被冻结
		ac.TokenNext(c, user) // 用户登录成功，生成 Token

		return
	} else {
		_ = global.Cache.Set(key, []byte(strconv.Itoa(count+1)))
		response.Response(c, http.StatusUnauthorized, fmt.Sprintf("用户 %v 登陆失败：验证码错误！", logonRequest.Username), nil)
		return
	}
}

func (ac *AuthController) TokenNext(c *gin.Context, user *model.User) {
	jwt, err := util.GenerateJWT(util.CustomClaims{
		UUID:        user.UUID,
		ID:          user.ID,
		Username:    user.Username,
		AuthorityId: user.AuthorityId,
	})
	if err != nil {
		global.Logger.Error(fmt.Sprintf("用户 %v 登陆失败：获取 token 失败！", user.Username), zap.Error(err))
		response.Response(c, http.StatusUnauthorized, fmt.Sprintf("用户 %v 登陆失败：获取 token 失败！", user.Username), nil)
	} else {
		global.Logger.Info(fmt.Sprintf("用户 %v 登陆成功！", user.Username))
		response.Response(c, http.StatusOK, "登陆成功！", gin.H{"token": jwt})
	}
}
