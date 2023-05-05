package middleware

import (
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/gogoclouds/gogo/logger"
	"github.com/gogoclouds/gogo/web/r"
	"io"
)

// LoggerResponseFail 日志记录中间件
// 1.只对 application/json
func LoggerResponseFail() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.ContentType() != binding.MIMEJSON { // 只输出 application/json
			c.Next()
			return
		}
		reqBody, _ := c.GetRawData()
		if len(reqBody) > 0 { // 请求包体写回。
			c.Request.Body = io.NopCloser(bytes.NewBuffer(reqBody))
		}
		blw := &responseWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
		c.Writer = blw
		c.Next()
		// 记录回包内容和处理时间
		rspBody := blw.body.Bytes()
		var rsp r.RespData[any]
		if err := json.Unmarshal(rspBody, &rsp); err != nil {
			logger.Error(err)
			return
		}
		if rsp.Code != r.Ok {
			logger.Errorf("\nrequest: %s\nresponse: %s", reqBody, rspBody)
		}
	}
}
