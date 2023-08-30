package insert

import (
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/joker-circus/gotools/internal"
)

// 生成批量创建 SQL 语句。
// 适合 Gorm v1 批量创建
func BuildInsertSQL(dest interface{}) (sql string, err error) {
	sqlBuild, err := buildSQL(dest)
	if err != nil {
		return "", err
	}

	return sqlBuild.ExplainSQL()
}

func BuildPreSQL(dest interface{}) (sql string, err error) {
	sqlBuild, err := buildSQL(dest)
	if err != nil {
		return "", err
	}

	return sqlBuild.PreSQL()
}

func buildSQL(dest interface{}) (SQL, error) {
	schema, err := GetSchema(dest)
	if err != nil {
		return SQL{}, err
	}

	now := reflect.ValueOf(time.Now())
	defaultValue := map[string]reflect.Value{
		"created_at": now,
		"updated_at": now,
	}

	// assign reflectValue
	reflectValue := reflect.ValueOf(dest)
	for reflectValue.Kind() == reflect.Ptr {
		if reflectValue.IsNil() && reflectValue.CanAddr() {
			reflectValue.Set(reflect.New(reflectValue.Type().Elem()))
		}

		reflectValue = reflectValue.Elem()
	}
	if !reflectValue.IsValid() {
		return SQL{}, ErrInvalidValue
	}

	return initSQL(schema, reflectValue, defaultValue)
}

func initSQL(schema *Schema, reflectValue reflect.Value, defaultValue map[string]reflect.Value) (SQL, error) {
	sqlBuild := SQL{
		Table:   schema.Table,
		Columns: schema.DBNames,
		Values:  nil,
	}

	// https://github.com/go-gorm/gorm/blob/v1.23.8/callbacks/create.go#L206
	switch reflectValue.Kind() {
	case reflect.Slice, reflect.Array:
		rValLen := reflectValue.Len()
		if rValLen == 0 {
			return sqlBuild, ErrEmptySlice
		}

		for i := 0; i < rValLen; i++ {
			rv := reflect.Indirect(reflectValue.Index(i))
			if !rv.IsValid() {
				return sqlBuild, fmt.Errorf("slice data #%v is invalid: %w", i, ErrUnsupportedDriver)
			}

			sqlBuild.Values = append(sqlBuild.Values, getStructValues(rv, schema, defaultValue))
		}

	case reflect.Struct:
		sqlBuild.Values = append(sqlBuild.Values, getStructValues(reflectValue, schema, defaultValue))

	default:
		return sqlBuild, ErrInvalidData
	}

	return sqlBuild, nil
}

// 获取结构体的所有值
func getStructValues(reflectValue reflect.Value, schema *Schema, defaultValue map[string]reflect.Value) []interface{} {
	values := make(map[string]interface{})
	r := &internal.StructX{T: reflectValue.Type(), V: reflectValue}
	r.RangeFields(false, func(sf reflect.StructField, rv reflect.Value) bool {
		field, ok := schema.FieldsByName[sf.Name]
		if !ok {
			return true
		}

		// 零值赋值
		if rv.IsZero() && rv.CanSet() {
			// sv.Type() == rv.Type()
			if sv, ok := defaultValue[field.DBName]; ok {
				values[field.DBName] = sv.Interface()
				return true
			}
		}

		// 对于软删除/指针 Time 做特殊处理
		// sf.Type.String() == "*time.Time"
		if sf.Type.Kind() == reflect.Ptr && rv.IsZero() && strings.HasSuffix(sf.Type.String(), "Time") {
			values[field.DBName] = nil
			return true
		}

		values[field.DBName] = rv.Interface()
		return true
	})

	data := make([]interface{}, len(schema.DBNames))
	for idx, column := range schema.DBNames {
		data[idx] = values[column]
	}
	return data
}
