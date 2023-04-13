package g

import (
	"github.com/gin-gonic/gin"
	"github.com/gogoclouds/gogo/web/r"
)

type IBaseApi interface {
	PageList(ctx *gin.Context)
	Create(ctx *gin.Context)
	Update(ctx *gin.Context)
	Delete(ctx *gin.Context)
}

type IBase[T any, Q any] interface {
	PageList(Q) (*r.PageResp[T], *Error)
	Create(T) *Error
	Update(T) *Error
	Delete(int) *Error
}