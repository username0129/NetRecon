package middleware

import (
	"backend/internal/global"
	"backend/internal/model/common"
	"backend/internal/service"
	"backend/internal/util"
	"github.com/gin-gonic/gin"
	"net/http"
)

var casbinService = service.CasbinServiceApp

func CasbinHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		if global.DB == nil {
			return
		}
		// 获取请求资源
		obj := c.Request.URL.Path
		// 获取请求方式
		act := c.Request.Method

		// 获取请求主体（身份 id）
		sub := util.GetAuthorityId(c)

		// 判断是否存在对应的 ACL
		casbin := casbinService.GetCasbin()

		ok, _ := casbin.Enforce(sub, obj, act)
		if ok {
			c.Next() // 请求成功
		} else {
			common.ResponseOk(c, http.StatusForbidden, "用户权限不足", nil)
			c.Abort() // 请求失败
			return
		}
	}
}
