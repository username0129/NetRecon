package controller

import (
	"backend/internal/global"
	"backend/internal/model/common"
	"backend/internal/model/request"
	"backend/internal/service"
	"backend/internal/util"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
)

type InitController struct{}

// GetInit
//
//	@Description: 检查数据库初始化状态
//	@receiver ic
//	@param c
//	@Router: /init/init
func (ic *InitController) GetInit(c *gin.Context) {
	if global.DB != nil {
		common.ResponseOk(c, http.StatusOK, "已存在数据库配置", nil)
		return
	}
	common.ResponseOk(c, http.StatusInternalServerError, "数据库尚未初始化", nil)
	return
}

// PostInit
//
//	@Description: 初始化数据库
//	@receiver ic
//	@param c
//	@Router: /init/init
func (ic *InitController) PostInit(c *gin.Context) {
	if global.DB != nil {
		global.Logger.Error("已存在数据库配置")
		common.ResponseOk(c, http.StatusInternalServerError, "已存在数据库配置", nil)
		return
	}

	var req request.InitRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		global.Logger.Error(fmt.Sprintf("参数解析错误 %v", err.Error()))
		common.ResponseOk(c, http.StatusInternalServerError, "参数解析错误", nil)
		return
	}

	if err := service.InitServiceApp.Init(req); err != nil {
		global.Logger.Error(fmt.Sprintf("数据库初始化错误：%v", err.Error()))
		common.ResponseOk(c, http.StatusInternalServerError, "数据库初始化错误，详情请查看后端。", nil)
		return
	}
	common.ResponseOk(c, http.StatusOK, "数据库初始化成功，请等待服务器重启... ", nil)

	administratorMail, err := service.UserServiceApp.GetAdministratorMail()
	if err != nil {
		global.Logger.Error("获取管理员邮箱失败: ", zap.Error(err))
	} else {
		body := `
<!DOCTYPE html>
<html>
<head>
  <style>
    body { font-family: 'Arial', sans-serif; line-height: 1.6; }
    h1 { color: #333; }
	p { margin: 10px 0; }
    .footer { color: grey; font-size: 0.9em; }
    hr { border: 0; height: 1px; background-color: #ddd; }
  </style>
</head>
<body>
  <h1>邮件系统测试</h1>
  <p>这是一封测试邮件，目的是验证邮件服务器配置是否成功。</p>
  <p>如果你能看到这封邮件，那么恭喜！你的邮件服务已经配置正确。</p>
  <hr>
  <p class="footer">请勿回复此邮件，此邮件为系统自动生成。</p>
</body>
</html>
`
		subject := "邮件服务配置测试"
		mail := global.Config.Mail
		err := util.SendMail(mail.SmtpServer, mail.SmtpPort, mail.SmtpFrom, mail.SmtpPassword, administratorMail, subject, body)
		if err != nil {
			global.Logger.Error("发送邮箱失败: ", zap.Error(err))
		}
	}
	// 触发服务器重启
	go func() {
		global.RestartSignal <- true
	}()

	return
}
