package types

import "reflect"

// 获取结构体中所有 tagName 值
func StructTagAllFields(dest interface{}, tagName string) (fields []string) {
	return StructTagFields(dest, tagName, true)
}

// 获取结构体中所有可导出字段 tagName 值
func StructTagExportedFields(dest interface{}, tagName string) (fields []string) {
	return StructTagFields(dest, tagName, false)
}

// StructTagFields 获取结构体中 tagName 值。
// isExported 表示是否获取不可导出字段值。
func StructTagFields(dest interface{}, tagName string, isExported bool) (fields []string) {
	rv := reflect.Indirect(reflect.ValueOf(dest))
	if rv.Kind() != reflect.Struct {
		return
	}

	return structTagFields(rv.Type(), tagName, isExported)
}

// 遍历 struct 所有字段的 tagName 值，
// rv.Kind() = reflect.Struct
func structTagFields(rvType reflect.Type, tagName string, isExported bool) (fields []string) {
	fields = make([]string, 0)
	for i := 0; i < rvType.NumField(); i++ {
		if rvType.Field(i).IsExported() || isExported {
			fields = append(fields, structValueTagFields(rvType.Field(i).Type, rvType.Field(i), tagName, isExported)...)
		}
	}
	return
}

// 获取 struct 某个字段的 tagName 值
func structValueTagFields(structType reflect.Type, structField reflect.StructField, tagName string, isExported bool) (fields []string) {
	// 去指针
	if structType.Kind() == reflect.Ptr {
		return structValueTagFields(structType.Elem(), structField, tagName, isExported)
	}

	// 嵌套结构体
	if structField.Anonymous && structType.Kind() == reflect.Struct {
		return structTagFields(structType, tagName, isExported)
	}

	tagValue := structField.Tag.Get(tagName)
	if len(tagValue) != 0 {
		fields = []string{tagValue}
	}

	return
}
