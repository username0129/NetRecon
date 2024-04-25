package common

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type ApiResponse struct {
	Code int         `json:"code"` // 响应码
	Msg  string      `json:"msg"`  // 响应消息
	Data interface{} `json:"data"` // 相应数据
}

func Response(c *gin.Context, code int, msg string, data interface{}) {
	c.JSON(http.StatusOK, ApiResponse{
		Code: code,
		Msg:  msg,
		Data: data,
	})
}
