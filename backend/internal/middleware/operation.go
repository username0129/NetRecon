package middleware

import (
	"backend/internal/global"
	"backend/internal/model"
	"backend/internal/util"
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid/v5"
	"go.uber.org/zap"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

// bodyWriter 重写 Response ，实现保存响应体数据
type bodyWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

// Write 将响应写入缓存和响应体
func (w *bodyWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

// parseBody 解析请求体
func parseBody(c *gin.Context) ([]byte, error) {
	// 使用 POST 方式
	if c.Request.Method != http.MethodGet {
		body, err := io.ReadAll(c.Request.Body)
		if err != nil {
			return nil, err
		}
		c.Request.Body = io.NopCloser(bytes.NewBuffer(body))
		return body, nil
	} else {
		// GET 方式
		query, _ := url.QueryUnescape(c.Request.URL.RawQuery)
		return parseQuery(query)
	}
}

// parseQuery 解析 URL 查询字符串
func parseQuery(query string) ([]byte, error) {
	params := make(map[string]string)
	for _, v := range strings.Split(query, "&") {
		kv := strings.SplitN(v, "=", 2)
		if len(kv) == 2 {
			params[kv[0]] = kv[1]
		}
	}
	return json.Marshal(params)
}

// OperationRecord 记录请求记录
func OperationRecord() gin.HandlerFunc {
	return func(c *gin.Context) {
		if global.DB == nil {
			return
		}

		body, err := parseBody(c) // 解析用户请求
		if err != nil {
			global.Logger.Error("解析请求体失败", zap.Error(err))
			return
		}
		writer := &bodyWriter{ResponseWriter: c.Writer, body: &bytes.Buffer{}}
		c.Writer = writer
		c.Next()

		resp := writer.body.String()

		// 检查请求体和响应体的长度
		if len(body) > 2048 {
			body = []byte(`{"error": "长度超出保存限制"}`)
		}
		if len(resp) > 2048 {
			resp = `{"error": "长度超出保存限制"}`
		}

		record := model.OperationRecord{
			UUID:     uuid.Must(uuid.NewV4()),
			IP:       c.ClientIP(),
			Method:   c.Request.Method,
			Path:     c.Request.URL.Path,
			Body:     string(body),
			UserUUID: util.GetUUID(c),
			Code:     strconv.Itoa(writer.Status()),
			Resp:     resp,
		}

		err = record.InsertData(global.DB)
		if err != nil {
			global.Logger.Error("记录操作失败: ", zap.Error(err))
		}
	}
}
