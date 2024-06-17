package router

import (
	"github.com/bobacgo/kit/app"
	"github.com/bobacgo/kit/examples/internal/app/admin"
	"github.com/gin-gonic/gin"
)

func Register(e *gin.Engine, app *app.Options) {
	admin.Register(e, app)
	// ....
}
