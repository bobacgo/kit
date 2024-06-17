package admin

import (
	"github.com/bobacgo/kit/app"
	"github.com/gin-gonic/gin"
)

func Register(e *gin.Engine, app *app.Options) {
	r := e.Group("")

	// sys user
	r.POST("v1/user/create", func(ctx *gin.Context) {})
	r.PUT("v1/user/update", func(ctx *gin.Context) {})
	r.DELETE("v1/user/delete", func(ctx *gin.Context) {})
	r.POST("v1/user/pageList", func(ctx *gin.Context) {})

	// ....
}
