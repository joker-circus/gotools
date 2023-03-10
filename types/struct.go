package types

import (
	"fmt"
	"reflect"
)

// 获取结构体中所有可导出 values 值，及 tagName 对应 structField 的映射关系
func StructTagExportedFieldValues(dest interface{}, tagName string) (tags []interface{}, tagFields map[string]string) {
	tagFields = make(map[string]string)
	r, ok := NewStructX(dest)
	if !ok {
		return
	}
	r.RangeFields(false, func(sf reflect.StructField, v reflect.Value) {
		tagValue := sf.Tag.Get(tagName)
		if len(tagValue) == 0 {
			return
		}

		if _, ok := tagFields[tagValue]; !ok {
			tagFields[tagValue] = sf.Name
			if v.CanInterface() {
				tags = append(tags, v.Interface())
			} else {
				tags = append(tags, fmt.Sprint(v))
			}
		}

	})
	return
}

// 获取结构体中所有 tagName 值，及 tagName 对应 structField 的映射关系
func StructTagAllFields(dest interface{}, tagName string) (tags []string, tagFields map[string]string) {
	return StructTagFields(dest, true, tagName)
}

// 获取结构体中所有可导出字段 tagName 值，及 tagName 对应 structField 的映射关系
func StructTagExportedFields(dest interface{}, tagName string) (tags []string, tagFields map[string]string) {
	return StructTagFields(dest, false, tagName)
}

// StructTagFields 获取结构体中 tagName 值。
// isExported 表示是否获取不可导出字段值。
func StructTagFields(dest interface{}, isExported bool, tagName string) (tags []string, tagFields map[string]string) {
	tagFields = make(map[string]string)
	r, ok := NewStructX(dest)
	if !ok {
		return
	}
	r.RangeFields(isExported, func(sf reflect.StructField, v reflect.Value) {
		tagValue := sf.Tag.Get(tagName)
		if len(tagValue) == 0 {
			return
		}

		if _, ok := tagFields[tagValue]; !ok {
			tagFields[tagValue] = sf.Name
			tags = append(tags, tagValue)
		}

	})
	return
}

type StructX struct {
	T reflect.Type
	V reflect.Value
}

// 返回 dest 结构体的 StructX 实例。
func NewStructX(dest interface{}) (*StructX, bool) {
	rv := reflect.Indirect(reflect.ValueOf(dest))
	if rv.Kind() != reflect.Struct {
		return nil, false
	}
	return &StructX{rv.Type(), rv}, true
}

// 返回 dest 结构体的 StructX 实例。
// 如果 dest 的类型不是 struct 类型或 struct 指针类型，结果会 panic。
func MustBeStructX(dest interface{}) *StructX {
	res, ok := NewStructX(dest)
	if !ok {
		panic("dest type must be struct")
	}
	return res
}

// isExported 表示是否遍历不可导出字段值。
func (r *StructX) RangeFields(isExported bool, f func(sf reflect.StructField, v reflect.Value)) {
	r.structTagFields(r.T, r.V, isExported, f)
}

// 遍历 struct 所有字段的 tagName 值，
// rv.Kind() = reflect.Struct
func (r *StructX) structTagFields(rvType reflect.Type, rv reflect.Value, isExported bool, f func(sf reflect.StructField, v reflect.Value)) {
	for i := 0; i < rvType.NumField(); i++ {
		if rvType.Field(i).IsExported() || isExported {
			r.structValueTagFields(rvType.Field(i).Type, rv.Field(i), rvType.Field(i), isExported, f)
		}
	}
	return
}

// 获取 struct 某个字段的 tagName 值
func (r *StructX) structValueTagFields(structType reflect.Type, structValue reflect.Value, structField reflect.StructField, isExported bool, f func(sf reflect.StructField, v reflect.Value)) {
	// 去指针
	if structType.Kind() == reflect.Ptr {
		if structValue.IsNil() {
			structValue = reflect.New(structType.Elem())
		}
		r.structValueTagFields(structType.Elem(), structValue.Elem(), structField, isExported, f)
		return
	}

	// 嵌套结构体
	if structField.Anonymous && structType.Kind() == reflect.Struct {
		r.structTagFields(structType, structValue, isExported, f)
		return
	}

	f(structField, structValue)
}
