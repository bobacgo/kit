package validator

import (
	"log/slog"

	"github.com/go-playground/locales"
	"github.com/go-playground/locales/en"
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	enTranslations "github.com/go-playground/validator/v10/translations/en"
	zhTranslations "github.com/go-playground/validator/v10/translations/zh"
	"golang.org/x/exp/maps"
)

// Trans 定义一个全局翻译器T
var trans *ut.UniversalTranslator

type TranslationLanguage struct {
	Lt           locales.Translator
	RegisterFunc func(*validator.Validate, ut.Translator) error
}

// AddTrans 添加翻译器
// 默认支持 en-英文和 zh-中文、zh_Hant_TW-繁体
// multipleTrans 支持其他国家或地区翻译器
func AddTrans(multipleTrans ...TranslationLanguage) {
	// 修改gin框架中的Validator引擎属性，实现自定制
	// 注册一个获取json tag的自定义方法
	registerTrans(validate, multipleTrans...)
}

func registerTrans(validate *validator.Validate, multipleTrans ...TranslationLanguage) {
	tMap := map[locales.Translator]func(*validator.Validate, ut.Translator) error{
		en.New(): enTranslations.RegisterDefaultTranslations,
		zh.New(): zhTranslations.RegisterDefaultTranslations,
	}
	for _, tran := range multipleTrans {
		tMap[tran.Lt] = tran.RegisterFunc
	}
	ts := maps.Keys(tMap)
	trans = ut.New(ts[0], ts...) // 默认 en
	for t, register := range tMap {
		lt, _ := trans.GetTranslator(t.Locale())
		if err := register(validate, lt); err != nil { // 注册多语言翻译器
			slog.Error("[validator] registerTrans", t.Locale(), err)
		}
	}
}
