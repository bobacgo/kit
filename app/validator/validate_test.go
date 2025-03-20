package validator_test

import (
	"context"
	"testing"

	"github.com/bobacgo/kit/app/validator"
	"github.com/go-playground/locales/ja"
	jaTranslations "github.com/go-playground/validator/v10/translations/ja"
)

// Address 用于测试嵌套结构体验证
type Address struct {
	Street  string `validate:"required"`
	City    string `validate:"required"`
	Country string `validate:"required,len=2"` // 国家代码必须是2个字符
}

// NestedExample 用于测试嵌套结构体验证
type NestedExample struct {
	Name    string   `validate:"required"`
	Address Address  `validate:"required"`                           // 嵌套结构体
	Tags    []string `validate:"required,min=1,dive,required,min=2"` // 切片验证
}

func TestJapaneseTranslation(t *testing.T) {
	// 注册日语翻译器
	validator.AddTrans(validator.TranslationLanguage{
		Lt:           ja.New(),
		RegisterFunc: jaTranslations.RegisterDefaultTranslations,
	})

	// 测试无效数据
	invalidData := NestedExample{
		Name:    "",            // 必填字段为空
		Address: Address{},     // 空的地址结构体
		Tags:    []string{"a"}, // 包含长度小于2的字符串
	}

	ctx := context.WithValue(context.Background(), "language", "zh")
	err := validator.StructCtx(ctx, invalidData)
	if err == nil {
		t.Error("Expected error for invalid data, got nil")
	}

	// 输出日语错误信息
	t.Logf("Japanese error message: %v", err)
}

func TestNestedStruct(t *testing.T) {
	// 测试有效数据
	validData := NestedExample{
		Name: "John Doe",
		Address: Address{
			Street:  "123 Main St",
			City:    "Beijing",
			Country: "CN",
		},
		Tags: []string{"tech", "golang", "testing"},
	}
	err := validator.Struct(validData)
	if err != nil {
		t.Errorf("Expected no error for valid nested data, got %v", err)
	}

	// 测试无效的嵌套结构体数据
	invalidData := NestedExample{
		Name:    "John Doe",
		Address: Address{}, // 空的地址结构体
		Tags:    []string{},
	}
	err = validator.Struct(invalidData)
	if err == nil {
		t.Error("Expected error for invalid nested data, got nil")
	}

	t.Log("invalidData:", err)

	// 测试无效的切片数据
	invalidSliceData := NestedExample{
		Name: "John Doe",
		Address: Address{
			Street:  "123 Main St",
			City:    "Beijing",
			Country: "CN",
		},
		Tags: []string{"a"}, // 包含长度小于2的字符串
	}
	err = validator.Struct(invalidSliceData)
	if err == nil {
		t.Error("Expected error for invalid slice data, got nil")
	}
}
