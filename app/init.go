package app

import (
	"context"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	zhTranslations "github.com/go-playground/validator/v10/translations/zh"
	"github.com/gogoclouds/gogo/g"
	"github.com/gogoclouds/gogo/internal/conf"
	"github.com/gogoclouds/gogo/internal/log"
	"reflect"
	"strings"
)

type app struct {
	ctx        context.Context
	conf       *conf.Config
	enableRpc  bool
	enableHttp bool
}

func Init() {
	initLogger()
	initTrans()
}

func initLogger() {
	log.Init("gogo", conf.Log{
		Level:       "debug", // debug | info | error
		FileSizeMax: 10,      // 10 MB
		FileAgeMax:  10,      // 10d
		DirPath:     "/logs",
	})
}

// initTrans 初始化翻译器
func initTrans() {
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
		g.Trans, _ = uni.GetTranslator("zh")
		if err := zhTranslations.RegisterDefaultTranslations(validate, g.Trans); err != nil {
			panic(err)
		}
	}
}