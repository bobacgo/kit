package middleware

import (
	"github.com/bobacgo/kit/web/r"
	"github.com/gin-gonic/gin"
)

func Recovery() func(c *gin.Context) {
	return gin.CustomRecovery(func(c *gin.Context, err any) {
		r.Reply(c, err)
	})
}
