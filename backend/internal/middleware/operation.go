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
	"time"
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

func OperationRecord() gin.HandlerFunc {
	return func(c *gin.Context) {
		if global.DB == nil {
			return
		}

		start := time.Now()

		var body []byte
		var err error
		// 如果使用其他方式如 POST，获取请求体中的内容
		if c.Request.Method != http.MethodGet {
			body, err = io.ReadAll(c.Request.Body)
			if err != nil {
				global.Logger.Error("读取请求体失败: ", zap.Error(err))
				return
			}
			// 重新设置读取指针
			c.Request.Body = io.NopCloser(bytes.NewBuffer(body))
		} else {
			// 如果使用 GET 方式，则获取 URL 中的内容
			query, _ := url.QueryUnescape(c.Request.URL.RawQuery)
			split := strings.Split(query, "&")
			m := make(map[string]string)
			for _, v := range split {
				kv := strings.SplitN(v, "=", 2)
				if len(kv) == 2 {
					m[kv[0]] = kv[1]
				}
			}
			body, err = json.Marshal(&m)
			if err != nil {
				global.Logger.Error("序列化查询参数失败", zap.Error(err))
				return
			}
		}

		writer := &bodyWriter{ResponseWriter: c.Writer, body: &bytes.Buffer{}}
		c.Writer = writer

		c.Next()

		duration := time.Since(start).Milliseconds() // 记录毫秒数

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
			Agent:    c.Request.UserAgent(),
			Body:     string(body),
			UserUUID: util.GetUUID(c),
			Duration: strconv.FormatInt(duration, 10),
			Code:     strconv.Itoa(writer.Status()),
			Resp:     resp,
		}

		if err := record.InsertData(global.DB); err != nil {
			global.Logger.Error("记录失败: ", zap.Error(err))
		}
	}
}
