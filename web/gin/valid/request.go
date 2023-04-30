package valid

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	zhTranslations "github.com/go-playground/validator/v10/translations/zh"
	"github.com/gogoclouds/gogo/logger"
	"github.com/gogoclouds/gogo/web/gin/reply"
	"github.com/gogoclouds/gogo/web/r"
	"reflect"
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
		if !ok { // 非validator.ValidationErrors类型错误直接返回
			logger.Error(err)
			reply.FailCodeDetails(ctx, r.ParameterIllegal, err.Error())
			return obj, false
		}
		// validator.ValidationErrors类型错误则进行翻译
		reply.FailCodeDetails(ctx, r.ParameterInvalid, removeTopStruct(errs.Translate(Trans)))
		return obj, false
	}
	return obj, true
}

// InitRequestParamValidate 初始化翻译器
func InitRequestParamValidate() {
	// 修改gin框架中的Validator引擎属性，实现自定制
	if validate, ok := binding.Validator.Engine().(*validator.Validate); ok {
		// 注册一个获取json tag的自定义方法
		validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
			name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
			// skip if tag key says it should be ignored
			if name == "-" {
				return ""
			}
			return name
		})

		uni := ut.New(zh.New())
		Trans, _ = uni.GetTranslator("zh")
		if err := zhTranslations.RegisterDefaultTranslations(validate, Trans); err != nil {
			panic(err)
		}
	}
}

// 去掉结构体名称前缀
func removeTopStruct(fields map[string]string) map[string]string {
	res := map[string]string{}
	for field, err := range fields {
		res[field[strings.LastIndex(field, ".")+1:]] = err
	}
	return res
}
