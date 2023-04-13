package valid

import (
	"github.com/gin-gonic/gin"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"github.com/gogoclouds/gogo/logger"
	"github.com/gogoclouds/gogo/web/gin/reply"
	"github.com/gogoclouds/gogo/web/r"
	"strings"
)

// Trans 定义一个全局翻译器T
var Trans ut.Translator

// ShouldBind 解析请求参数并校验
// T 返回 请求数据
// 解析失败或者参数校验失败，返回 false  并组织响应数据结构
// 参数校验成功，返回 true
//
// req, ok := valid.ShouldBind[model.CreateUserReq](ctx)
//
//	if ok {
//	    return
//	}
func ShouldBind[T any](ctx *gin.Context) (T, bool) {
	var obj T
	if err := ctx.ShouldBind(&obj); err != nil {
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
		res[field[strings.LastIndex(field, ".")+1:]] = err
	}
	return res
}