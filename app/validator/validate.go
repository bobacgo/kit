package validator

import (
	"context"
	"errors"
	"strings"

	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

func init() {
	// 注册翻译器 默认支持 zh en
	// AddTrans() 添加语种
	registerTrans(validate)
	// 注册自定义验证器
	registerValidation(validate)
}

func Get() *validator.Validate {
	return validate
}

// Validate 验证器
func Struct(obj any) error {
	if err := validate.Struct(obj); err != nil {
		return TransErr(err)
	}
	return nil
}

// StructCtx 验证器
// 支持 context ctx.Value(LanguageCtxKey)
func StructCtx(ctx context.Context, obj any) error {
	if err := validate.StructCtx(ctx, obj); err != nil {
		return TransErrCtx(ctx, err)
	}
	return nil
}

// TransErrZh 解析错误信息为中文
func TransErrZh(err error) error {
	return TransErrLocale("zh", err)
}

// TransErr 解析错误信息为英文
func TransErr(err error) error {
	return TransErrLocale("en", err)
}

// TransErrCtx 解析错误信息 支持 context ctx.Value(LanguageCtxKey)
func TransErrCtx(ctx context.Context, err error) error {
	return TransErrLocale(DefaultGetLanguage(ctx), err)
}

// TransErrLocale 解析错误信息
// 支持 指定语言
func TransErrLocale(locale string, err error) error {
	verr, ok := err.(validator.ValidationErrors)
	if !ok {
		return err
	}

	t, _ := trans.GetTranslator(locale)
	msgErr := removeTopStruct(verr.Translate(t))
	return errors.New(msgErr)
}

// 去掉最顶层结构体名称前缀，但保留嵌套结构体的路径
func removeTopStruct(fields map[string]string) string {
	msgErrs := strings.Builder{}
	for field, err := range fields {
		// 分割字段路径
		parts := strings.Split(field, ".")
		// 如果有嵌套结构体，保留除了最顶层以外的路径
		if len(parts) > 1 {
			field = strings.Join(parts[1:], ".")
			msgErrs.WriteString(field)
			msgErrs.WriteString(": ")
		}
		msgErrs.WriteString(err)
		msgErrs.WriteString(", ")
	}
	return strings.TrimRight(msgErrs.String(), ", ")
}
