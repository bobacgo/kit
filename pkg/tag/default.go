package tag

import (
	"fmt"
	"golang.org/x/exp/slices"
	"log/slog"
	"reflect"
	"strconv"
)

// Default 处理默认赋值的标签
func Default[T any](data T) T {
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
	new(defaultTag).parseStruct(val, newVal)

	// 返回新的对象
	if isPtr {
		return newVal.Addr().Interface().(T)
	}
	return newVal.Interface().(T)
}

type defaultTag struct{}

// parseStruct 只处理 struct、ptr、map、slice、array
func (t *defaultTag) parseStruct(src, dst reflect.Value) {
	switch src.Kind() {
	case reflect.Struct: // 处理结构体
		for i := 0; i < src.NumField(); i++ {
			field := src.Field(i)
			structField := src.Type().Field(i)
			newField := reflect.New(field.Type()).Elem()

			if !slices.Contains(basicKind, src.Kind()) {
				t.parseStruct(field, newField)
			} else {
				t.tagValue(field, newField, structField)
			}
			dst.Field(i).Set(newField)
		}
	case reflect.Ptr: // 处理接口类型和指针类型
		if src.IsNil() {
			return
		}
		elem := src.Elem()
		newElement := reflect.New(elem.Type()).Elem()
		t.parseStruct(elem, newElement)
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
			t.parseStruct(element, newElement)
			newContainer.Index(i).Set(newElement)
		}
		dst.Set(newContainer)
	// 创建一个新的 Slice 或 Array
	case reflect.Map:
		if src.IsNil() {
			return
		}
		// 创建一个新的 Map
		newMap := reflect.MakeMap(src.Type())
		for _, key := range src.MapKeys() {
			// 递归处理 Map 的每个值
			value := src.MapIndex(key)
			newValue := reflect.New(value.Type()).Elem()
			t.parseStruct(value, newValue)
			newMap.SetMapIndex(key, newValue)
		}
		dst.Set(newMap)
	default:
		dst.Set(src)
	}
}

func (t *defaultTag) tagValue(src, dst reflect.Value, structField reflect.StructField) {
	// 检查是否有 default 标签
	tagValue, ok := structField.Tag.Lookup("default")
	if !ok || tagValue == "" {
		dst.Set(src)
		return
	}
	if !src.IsZero() { // 如果源值不为零，则不设置默认值
		return
	}
	if err := t.set(dst, src.Kind(), tagValue); err != nil {
		slog.Warn("default tag value parse error", "err", err)
	}
}

func (t *defaultTag) set(dst reflect.Value, kind reflect.Kind, val string) error {
	switch kind {
	case reflect.Bool:
		dstVal, err := strconv.ParseBool(val)
		if err != nil {
			return fmt.Errorf("to Bool: %v", err)
		}
		dst.SetBool(dstVal)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		dstVal, err := strconv.ParseInt(val, 10, dst.Type().Bits())
		if err != nil {
			return fmt.Errorf("to Int: %v", err)
		}
		dst.SetInt(dstVal)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		dstVal, err := strconv.ParseUint(val, 10, dst.Type().Bits())
		if err != nil {
			return fmt.Errorf("to Uint: %v", err)
		}
		dst.SetUint(dstVal)
	case reflect.String:
		dst.SetString(val)
	default:
	}
	return nil
}