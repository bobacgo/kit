package g

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/gogoclouds/gogo/logger"
	"github.com/gogoclouds/gogo/web/gin/reply"
	"github.com/gogoclouds/gogo/web/r"
	"strings"
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

func BindAndValidate[T any](ctx *gin.Context) (T, bool) {
	var obj T
	if err := ctx.ShouldBind(&obj); err != nil {
		// 获取validator.ValidationErrors类型的errors
		errs, ok := err.(validator.ValidationErrors)
		if !ok {
			// 非validator.ValidationErrors类型错误直接返回
			logger.Error(err)
			reply.FailMsg(ctx, r.FailCreate)
			return obj, false
		}
		// validator.ValidationErrors类型错误则进行翻译
		reply.FailMsgDetails(ctx, r.FailCreate, removeTopStruct(errs.Translate(Trans)))
		return obj, false
	}
	return obj, true
}

// 去掉结构体名称前缀
func removeTopStruct(fields map[string]string) map[string]string {
	res := map[string]string{}
	for field, err := range fields {
		res[field[strings.Index(field, ".")+1:]] = err
	}
	return res
}