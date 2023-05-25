package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/gogoclouds/gogo/web/gin/reply"
	"github.com/gogoclouds/gogo/web/r"
)

func Recovery() func(c *gin.Context) {
	return gin.CustomRecovery(func(c *gin.Context, err any) {
		reply.FailCode(c, r.Internal)
	})
}
