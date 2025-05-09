package g

import (
	"github.com/bobacgo/kit/web/r/page"
	"github.com/gin-gonic/gin"
)

type IBaseApi interface {
	PageList(ctx *gin.Context)
	Create(ctx *gin.Context)
	Update(ctx *gin.Context)
	Delete(ctx *gin.Context)
}

type IBase[T any, Q any] interface {
	PageList(Q) (*page.Data[T], error)
	Create(T) error
	Update(T) error
	Delete(int) error
}
