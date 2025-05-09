package tag

import (
	"log/slog"
	"reflect"
	"regexp"
	"slices"
	"strings"
)

// Desensitize 对结构体中的敏感字段进行脱敏处理
// 只支持字符串类型
// e.g.
// Email    string   `json:"email" mask:""`        // 使用默认规则
// Password string   `json:"password" mask:"^.*$"` // 使用正则表达式
func Desensitize[T any](data T) T {
	val := reflect.ValueOf(data)
	if val.IsZero() {
		return data
	}
	isPtr := val.Kind() == reflect.Ptr
	if isPtr {
		val = val.Elem()
	}
	if val.Kind() != reflect.Struct { // 如果不是结构体，直接返回
		return data
	}
	// 创建一个新的同类型对象
	newVal := reflect.New(val.Type()).Elem()

	// 递归处理结构体字段
	new(maskTag).parseStruct(val, newVal, reflect.StructField{})

	// 返回新的对象
	if isPtr {
		return newVal.Addr().Interface().(T)
	}
	return newVal.Interface().(T)
}

type maskTag struct{}

// parseStruct 只处理 struct、ptr、map、slice、array
func (t *maskTag) parseStruct(src, dst reflect.Value, structField reflect.StructField) {
	switch src.Kind() {
	case reflect.Struct: // 处理结构体
		for i := 0; i < src.NumField(); i++ {
			field := src.Field(i)
			structField = src.Type().Field(i)
			newField := reflect.New(field.Type()).Elem()
			t.set(field, newField, structField)
			dst.Field(i).Set(newField)
		}
	case reflect.Ptr: // 处理接口类型和指针类型
		if src.IsNil() {
			return
		}
		elem := src.Elem()
		newElement := reflect.New(elem.Type()).Elem()
		t.set(elem, newElement, structField)
		dst.Set(newElement.Addr())
	case reflect.Slice, reflect.Array: // 处理 Slice 和 Array
		if src.Kind() == reflect.Slice && src.IsNil() {
			return
		}
		// 创建一个新的 Slice 或 Array
		var newContainer reflect.Value
		if src.Kind() == reflect.Slice {
			newContainer = reflect.MakeSlice(src.Type(), src.Len(), src.Cap())
		} else {
			newContainer = reflect.New(src.Type()).Elem()
		}

		for i := 0; i < src.Len(); i++ {
			// 递归处理 Slice 或 Array 的每个元素
			element := src.Index(i)
			newElement := reflect.New(element.Type()).Elem()
			t.set(element, newElement, structField)
			newContainer.Index(i).Set(newElement)
		}
		dst.Set(newContainer)
	case reflect.Map: // 处理 Map
		if src.IsNil() {
			return
		}
		// 创建一个新的 Map
		newMap := reflect.MakeMap(src.Type())
		for _, key := range src.MapKeys() {
			// 递归处理 Map 的每个值
			value := src.MapIndex(key)
			newValue := reflect.New(value.Type()).Elem()
			t.set(value, newValue, structField)
			newMap.SetMapIndex(key, newValue)
		}
		dst.Set(newMap)
	default: // 其他类型直接赋值
		dst.Set(src)
	}
}

func (t *maskTag) set(field, newField reflect.Value, fieldType reflect.StructField) {
	// 如果字段是Struct、Slice、Array、Map 类型，递归处理
	if slices.Contains(multiElement, field.Kind()) {
		t.parseStruct(field, newField, fieldType)
		return
	}
	if field.Kind() == reflect.String {
		// 检查是否有 mask 标签
		if maskTag, ok := fieldType.Tag.Lookup("mask"); ok {
			// 执行脱敏逻辑
			maskedValue := t.maskString(field.String(), maskTag)
			newField.SetString(maskedValue)
			return
		}
	}
	newField.Set(field)
}

// maskString 对字符串进行脱敏
func (t *maskTag) maskString(value, maskTag string) string {
	if maskTag == "" {
		// 使用默认规则：替换中间三分之一部分
		return t.defaultMask(value)
	}

	// 使用正则表达式进行脱敏
	re, err := regexp.Compile(maskTag)
	if err != nil {
		slog.Error("Error compiling regex:", maskTag, err)
		return value // 如果正则表达式错误，保留原始值
	}

	return re.ReplaceAllStringFunc(value, func(match string) string {
		return strings.Repeat("*", len(match))
	})
}

// 默认脱敏规则：替换中间三分之一部分
func (t *maskTag) defaultMask(value string) string {
	length := len(value)
	if length == 0 {
		return value
	}

	// 计算需要替换的起始和结束位置
	start := length / 3
	end := 2 * length / 3

	// 替换中间三分之一部分
	masked := value[:start] + strings.Repeat("*", end-start) + value[end:]
	return masked
}