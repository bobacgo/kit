package middleware

import (
	"errors"

	"github.com/bobacgo/kit/web/r"
	"github.com/gin-gonic/gin"
)

func Recovery() func(c *gin.Context) {
	return gin.CustomRecovery(func(c *gin.Context, err any) {
		if errMsg, ok := err.(string); ok {
			err = errors.New(errMsg)
		}
		r.Reply(c, err)
	})
}
