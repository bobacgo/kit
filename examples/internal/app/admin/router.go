package admin

import (
	"github.com/bobacgo/kit/app"
	"github.com/bobacgo/kit/examples/internal/app/admin/handler"
	"github.com/gin-gonic/gin"
)

func Register(r *gin.RouterGroup, app *app.AppOptions) {

	userHandler := handler.NewUserHandler()
	// sys user
	r.GET("v1/user/:id", userHandler.Get)
	r.POST("v1/user/create", userHandler.Create)
	r.PUT("v1/user/update", userHandler.Update)
	r.DELETE("v1/user/delete", userHandler.Delete)
	r.POST("v1/user/pageList", userHandler.PageList)

	// ....
}
