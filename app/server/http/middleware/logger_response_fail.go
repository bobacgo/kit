package middleware

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"

	"github.com/bobacgo/kit/web/r"
	"github.com/bobacgo/kit/web/r/codes"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
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
		var rsp r.Response[any]
		if err := json.Unmarshal(rspBody, &rsp); err != nil {
			return
		}
		if rsp.Code != codes.OK {
			slog.ErrorContext(c, fmt.Sprintf("\nrequest: %s\nresponse: %s", reqBody, rspBody))
		}
	}
}
