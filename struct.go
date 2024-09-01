package gotools

import (
	"fmt"
	"reflect"

	"github.com/joker-circus/gotools/internal"
)

// 获取结构体中所有可导出 values 值，及 tagFieldName 对应 structFieldName 的映射关系
func StructTagExportedFieldValues(dest interface{}, tagName string) (fieldValues []interface{}, tagFields map[string]string) {
	tagFields = make(map[string]string)
	r, ok := NewStructX(dest)
	if !ok {
		return
	}
	r.RangeFields(false, func(sf reflect.StructField, v reflect.Value) bool {
		fieldName, _ := StructTagFieldName(sf.Tag, tagName)
		if len(fieldName) == 0 {
			return true
		}

		if _, ok := tagFields[fieldName]; !ok {
			tagFields[fieldName] = sf.Name
			if v.CanInterface() {
				fieldValues = append(fieldValues, v.Interface())
			} else {
				fieldValues = append(fieldValues, fmt.Sprint(v))
			}
		}
		return true
	})
	return
}

// 获取结构体中所有 tagFieldName 值，及 tagFieldName 对应 structFieldName 的映射关系
func StructTagAllFields(dest interface{}, tagName string) (fields []string, tagFields map[string]string) {
	return StructTagFields(dest, true, tagName)
}

// 获取结构体中所有可导出字段 tagFieldName 值，及 tagFieldName 对应 structFieldName 的映射关系
func StructTagExportedFields(dest interface{}, tagName string) (fields []string, tagFields map[string]string) {
	return StructTagFields(dest, false, tagName)
}

// StructTagFields 获取结构体中 tagFieldName 值。
// isExported 表示是否获取不可导出字段值。
func StructTagFields(dest interface{}, isExported bool, tagName string) (fields []string, tagFields map[string]string) {
	tagFields = make(map[string]string)
	r, ok := NewStructX(dest)
	if !ok {
		return
	}
	r.RangeFields(isExported, func(sf reflect.StructField, v reflect.Value) bool {
		fieldName, _ := StructTagFieldName(sf.Tag, tagName)
		if len(fieldName) == 0 {
			return true
		}

		if _, ok := tagFields[fieldName]; !ok {
			tagFields[fieldName] = sf.Name
			fields = append(fields, fieldName)
		}

		return true
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

// 结构体数组转 Table 数据
func StructArrayToTable(dest interface{}, tagNames ...string) (columns []string, rows [][]interface{}) {
	rv := reflect.Indirect(reflect.ValueOf(dest))
	if rv.Kind() != reflect.Slice {
		return
	}

	if rv.Len() == 0 {
		return
	}

	// json 标签兜底
	tagNames = append(tagNames, "json")

	for i := 0; i < rv.Len(); i++ {
		r := StructX{
			T: rv.Index(i).Type(), // 必须逐步获取元素类型
			V: rv.Index(i),
		}
		if rv.Index(i).Kind() == reflect.Ptr {
			r.V = r.V.Elem()
			r.T = r.T.Elem()
		}
		if r.T.Kind() != reflect.Struct {
			continue
		}
		tagFields := make(map[string]string)
		var tags []interface{}
		r.RangeFields(false, func(sf reflect.StructField, v reflect.Value) bool {
			fieldName, _ := StructFieldNameByTagNames(sf, tagNames...)
			if len(fieldName) == 0 {
				return true
			}

			if i == 0 {
				columns = append(columns, fieldName)
			}

			if _, ok := tagFields[fieldName]; !ok {
				tagFields[fieldName] = sf.Name
				if v.CanInterface() {
					tags = append(tags, v.Interface())
				} else {
					tags = append(tags, fmt.Sprint(v))
				}
			}
			return true
		})
		rows = append(rows, tags)
	}
	return
}

// 根据 dest 结构体，判断 requiredFields 是否都已赋值。
func ValidateStruct(dest interface{}, tagName string, requiredFields ...string) (err error) {
	if len(requiredFields) == 0 {
		return nil
	}

	r, ok := NewStructX(dest)
	if !ok {
		return nil
	}

	requiredMap := make(map[string]struct{}, len(requiredFields))
	for _, s := range requiredFields {
		requiredMap[s] = struct{}{}
	}

	r.RangeFields(false, func(sf reflect.StructField, v reflect.Value) bool {
		fieldName, ok := StructTagFieldName(sf.Tag, tagName)
		if !ok {
			return true
		}

		_, ok = requiredMap[fieldName]
		if !ok {
			return true
		}

		if v.IsValid() && v.IsZero() {
			err = fmt.Errorf("%s 为必填字段", fieldName)
			return false
		}
		return true
	})

	return
}

// 根据 dest 结构体映射对应的值。
// 若 whitelist 有值，则仅获取 whitelist 内的字段值。
func StructToMap(dest interface{}, tagName string, whitelist ...string) (fields map[string]interface{}, err error) {
	r, ok := NewStructX(dest)
	if !ok {
		return nil, fmt.Errorf("dest must be struct")
	}

	whitelistMap := make(map[string]struct{})
	for _, v := range whitelist {
		whitelistMap[v] = struct{}{}
	}
	filter := func(fieldName string) (exist bool) {
		exist = true
		if len(whitelistMap) > 0 {
			_, exist = whitelistMap[fieldName]
		}
		return exist
	}

	fields = make(map[string]interface{})
	r.RangeFields(false, func(sf reflect.StructField, v reflect.Value) bool {
		fieldName, ok := StructTagFieldName(sf.Tag, tagName)
		if !ok {
			return true
		}

		// 过滤不需要的字段
		if !filter(fieldName) {
			return true
		}

		if v.CanInterface() {
			fields[fieldName] = v.Interface()
		} else {
			fields[fieldName] = fmt.Sprint(v)
		}

		return true
	})
	return
}

// 获取结构体字段的名称。
// 若 tagNames 不为空则，则获取 tagName 优先获取到的值，
// 否则默认获取 structFieldName 对应的蛇形/下划线字段名（例如：xx_y_y）。
func StructFieldName(sf reflect.StructField, tagNames ...string) string {
	fieldName, ok := StructFieldNameByTagNames(sf, tagNames...)
	if ok {
		return fieldName
	}
	return SnakeCase(sf.Name)
}

// 获取结构体字段的名称
// 例如 json:"name,omitempty"、gorm:"column:name" 返回的都是 name
func StructFieldNameByTagNames(sf reflect.StructField, tagNames ...string) (fieldName string, ok bool) {
	for _, tagName := range tagNames {
		if fieldName, ok = StructTagFieldName(sf.Tag, tagName); ok && fieldName != "" {
			return fieldName, true
		}
	}
	return "", false
}

// 获取结构体字段 Tag 对应的 FieldName 值，
// 例如 json:"name,omitempty"、gorm:"column:name" 返回的都是 name
func StructTagFieldName(tag reflect.StructTag, tagName string) (value string, ok bool) {
	fieldName, ok := tag.Lookup(tagName)
	if tagName == "gorm" {
		fieldName, ok = internal.GetGormTagColumnName(fieldName)
	}
	return fieldName, ok
}
