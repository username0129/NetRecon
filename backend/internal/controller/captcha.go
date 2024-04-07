package controller

import (
	"backend/internal/model"
	"backend/internal/model/response"
	"errors"
	"github.com/allegro/bigcache/v3"
	"github.com/gin-gonic/gin"
	"github.com/mojocn/base64Captcha"
	"go.uber.org/zap"
	"net/http"
	"strconv"

	"backend/internal/global"
)

type CaptchaController struct{}

func (cc *CaptchaController) GetCaptcha(c *gin.Context) {
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

	var oc bool
	if openCaptcha == 0 || openCaptcha <= count {
		oc = true
	}

	driver := base64Captcha.NewDriverDigit(global.Config.Captcha.ImgHeight, global.Config.Captcha.ImgWidth, global.Config.Captcha.Long, 0.7, 80)
	cp := base64Captcha.NewCaptcha(driver, global.CaptchaStore)
	id, b64s, _, err := cp.Generate()
	if err != nil {
		global.Logger.Error("验证码获取失败", zap.Error(err))
		response.Response(c, http.StatusInternalServerError, "验证码获取失败", nil)
		return
	}

	response.Response(c, http.StatusOK, "验证码获取成功", model.CaptchaResponse{
		CaptchaId:     id,
		CaptchaImg:    b64s,
		CaptchaLength: global.Config.Captcha.Long,
		OpenCaptcha:   oc,
	})
}
