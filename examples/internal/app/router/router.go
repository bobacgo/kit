package router

import (
	"github.com/bobacgo/kit/app"
	"github.com/bobacgo/kit/examples/internal/app/admin"
	"github.com/bobacgo/kit/examples/internal/app/middleware"
	"github.com/gin-gonic/gin"
)

func Register(e *gin.Engine, app *app.AppOptions) {
	r := e.Group("/")
	r.Use(middleware.Auth()) // 使用鉴权中间件
	admin.Register(r, app)
	// ....
}
