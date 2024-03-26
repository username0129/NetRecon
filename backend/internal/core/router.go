package core

import (
	"backend/internal/controller"
	"backend/internal/global"
	"backend/internal/middleware"
	"github.com/gin-gonic/gin"
	"reflect"
	"strings"
)

func InitializeRout() *gin.Engine {
	r := gin.New()
	r.Use(gin.Recovery()) // 避免 panic 导致服务器停止

	_ = r.SetTrustedProxies(nil) // 设置信任网络 nil 为不计算，避免性能消耗
	setupRoutes(r, getControllerList())

	return r
}

func getControllerList() []interface{} {
	return []interface{}{
		&controller.BaseController{},
		&controller.UserController{},
		&controller.AuthController{},
		&controller.CasbinController{},
		&controller.InitController{},
		&controller.CaptchaController{},
	}
}

func setupRoutes(router *gin.Engine, controllers []interface{}) {
	publicGroup := router.Group(global.Config.System.RouterPrefix) // 无需鉴权的路由组
	publicGroup.Use(middleware.CorsMiddleware())

	{
		// 全局处理 OPTIONS 请求 -> 防止 AXIOS CORS 错误
		publicGroup.OPTIONS("/*any", func(c *gin.Context) {
			// 可以选择返回一些内容或状态码
			c.JSON(200, gin.H{"message": "ok"})
		})

	}
	protectedGroup := router.Group(global.Config.System.RouterPrefix)                                                   // 需要鉴权的路由组
	protectedGroup.Use(middleware.CorsMiddleware()).Use(middleware.JWTAuthMiddleware()).Use(middleware.CasbinHandler()) // 使用 JWT 和 Casbin 完成身份验证以及访问控制

	for _, ctrl := range controllers {
		ctrlType := reflect.TypeOf(ctrl)
		ctrlValue := reflect.ValueOf(ctrl)
		ctrlName := strings.TrimSuffix(ctrlType.Elem().Name(), "Controller")

		_, jwtRequired := ctrlType.Elem().FieldByName("JWTRequired") // 检查是否需要鉴权

		for i := 0; i < ctrlType.NumMethod(); i++ {
			method := ctrlType.Method(i)
			methodName := method.Name

			if httpMethod, ok := getHTTPMethodFromName(methodName); ok {
				path := "/" + strings.ToLower(ctrlName) + "/" + strings.ToLower(strings.TrimPrefix(methodName, httpMethod)) // 创建路径

				handleFunc := func(c *gin.Context) {
					ctrlValue.MethodByName(methodName).Call([]reflect.Value{reflect.ValueOf(c)})
				} // 通过反射调用对应的方法

				if jwtRequired {
					protectedGroup.Handle(httpMethod, path, handleFunc)
				} else {
					publicGroup.Handle(httpMethod, path, handleFunc)
				}
			}
		} // 遍历 Controller 中实现的方法并添加到路由组
	}
}

func getHTTPMethodFromName(methodName string) (string, bool) {
	// 使用前缀匹配方法名，确定对应的HTTP方法
	if strings.HasPrefix(methodName, "Get") {
		return "GET", true
	} else if strings.HasPrefix(methodName, "Post") {
		return "POST", true
	} else if strings.HasPrefix(methodName, "Put") {
		return "PUT", true
	} else if strings.HasPrefix(methodName, "Delete") {
		return "DELETE", true
	} else if strings.HasPrefix(methodName, "Patch") {
		return "PATCH", true
	}
	// 如果没有匹配到任何HTTP方法，则返回false
	return "", false
}
