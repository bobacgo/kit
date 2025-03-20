package validator

import (
	"context"
)

var languageCtxKey = "language"

// SetLanguageCtxKey 设置语言上下文键值
func SetLanguageCtxKey(key string) {
	languageCtxKey = key
}

func GetLanguageCtxKey() string {
	return languageCtxKey
}

type GetRequestLanguageFunc func(ctx context.Context) string

// DefaultGetLanguage 重新赋值修改默认获取语言的方法
var DefaultGetLanguage GetRequestLanguageFunc = func(ctx context.Context) string {
	lang, _ := ctx.Value(languageCtxKey).(string)
	if lang == "" {
		lang = "en"
	}
	return lang
}
