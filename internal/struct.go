package internal

import (
	"reflect"
)

type StructX struct {
	T reflect.Type
	V reflect.Value
}

// isExported 表示是否遍历不可导出字段值。
func (r *StructX) RangeFields(isExported bool, f func(sf reflect.StructField, v reflect.Value) bool) {
	r.structTagFields(r.T, r.V, isExported, f)
}

// 遍历 struct 所有字段的 tagName 值，
// rv.Kind() = reflect.Struct
func (r *StructX) structTagFields(rvType reflect.Type, rv reflect.Value, isExported bool, f func(sf reflect.StructField, v reflect.Value) bool) bool {
	for i := 0; i < rvType.NumField(); i++ {
		if rvType.Field(i).IsExported() || isExported {
			if !r.structValueTagFields(rvType.Field(i).Type, rv.Field(i), rvType.Field(i), isExported, f) {
				return false
			}
		}
	}
	return true
}

// 获取 struct 某个字段的 tagName 值
func (r *StructX) structValueTagFields(structType reflect.Type, structValue reflect.Value, structField reflect.StructField, isExported bool, f func(sf reflect.StructField, v reflect.Value) bool) bool {
	// 去指针
	if structType.Kind() == reflect.Ptr {
		if structValue.IsNil() {
			structValue = reflect.New(structType.Elem())
		}
		return r.structValueTagFields(structType.Elem(), structValue.Elem(), structField, isExported, f)
	}

	// 嵌套结构体
	if structField.Anonymous && structType.Kind() == reflect.Struct {
		return r.structTagFields(structType, structValue, isExported, f)
	}

	return f(structField, structValue)
}
