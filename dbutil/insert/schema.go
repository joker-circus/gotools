package insert

import (
	"database/sql/driver"
	"fmt"
	"go/ast"
	"reflect"
)

type Schema struct {
	Name      string
	ModelType reflect.Type
	Table     string

	DBNames        []string
	Fields         []*Field
	FieldsByName   map[string]*Field
	FieldsByDBName map[string]*Field
}

type Field struct {
	Name           string
	DBName         string
	TagSettings    map[string]string
	Schema         *Schema
	EmbeddedSchema *Schema
}

// https://github.com/go-gorm/gorm/blob/v1.23.8/schema/schema.go#L80
func GetSchema(dest interface{}) (*Schema, error) {
	if dest == nil {
		return nil, fmt.Errorf("%w: %+v", ErrUnsupportedDataType, dest)
	}

	value := reflect.ValueOf(dest)
	if value.Kind() == reflect.Ptr && value.IsNil() {
		value = reflect.New(value.Type().Elem())
	}
	modelType := reflect.Indirect(value).Type()

	if modelType.Kind() == reflect.Interface {
		modelType = reflect.Indirect(reflect.ValueOf(dest)).Elem().Type()
	}

	for modelType.Kind() == reflect.Slice || modelType.Kind() == reflect.Array || modelType.Kind() == reflect.Ptr {
		modelType = modelType.Elem()
	}

	if modelType.Kind() != reflect.Struct {
		if modelType.PkgPath() == "" {
			return nil, fmt.Errorf("%w: %+v", ErrUnsupportedDataType, dest)
		}
		return nil, fmt.Errorf("%w: %s.%s", ErrUnsupportedDataType, modelType.PkgPath(), modelType.Name())
	}

	modelValue := reflect.New(modelType)
	tableName := modelValue.String()
	if tabler, ok := modelValue.Interface().(Tabler); ok {
		tableName = tabler.TableName()
	}

	schema := &Schema{
		Name:           modelType.Name(),
		ModelType:      modelType,
		Table:          tableName,
		FieldsByName:   map[string]*Field{},
		FieldsByDBName: map[string]*Field{},
	}

	var err error
	for i := 0; i < modelType.NumField(); i++ {
		if fieldStruct := modelType.Field(i); ast.IsExported(fieldStruct.Name) {
			tagSetting := ParseTagSetting(fieldStruct.Tag.Get("gorm"), ";")

			name := fieldStruct.Name
			dbName := tagSetting["COLUMN"]
			if dbName == "" {
				dbName = toDBName(name)
			}

			field := &Field{
				Name:        name,
				DBName:      dbName,
				TagSettings: tagSetting,
				Schema:      schema,
			}

			field.EmbeddedSchema, err = schema.embeddedSchema(fieldStruct)
			if err != nil {
				return nil, err
			}

			if field.EmbeddedSchema != nil {
				schema.Fields = append(schema.Fields, field.EmbeddedSchema.Fields...)
				continue
			}

			schema.Fields = append(schema.Fields, field)
		}
	}

	for i := range schema.Fields {
		field := schema.Fields[i]
		schema.FieldsByDBName[field.DBName] = field
		schema.FieldsByName[field.Name] = field
		schema.DBNames = append(schema.DBNames, field.DBName)
	}

	return schema, nil
}

// 是否为内嵌结构体
func (schema *Schema) embeddedSchema(fieldStruct reflect.StructField) (*Schema, error) {
	fieldType := fieldStruct.Type
	for fieldType.Kind() == reflect.Ptr {
		fieldType = fieldType.Elem()
	}

	fieldValue := reflect.New(fieldType)

	// if field is valuer, used its value or first field as data type
	valuer, isValuer := fieldValue.Interface().(driver.Valuer)
	if isValuer {
		if _, ok := fieldValue.Interface().(GormDataTypeInterface); !ok {
			if v, err := valuer.Value(); reflect.ValueOf(v).IsValid() && err == nil {
				fieldValue = reflect.ValueOf(v)
			}
		}
	}

	if fieldStruct.Anonymous {
		switch reflect.Indirect(fieldValue).Kind() {
		case reflect.Struct:
			return getOrParse(fieldValue.Interface())
		case reflect.Invalid, reflect.Uintptr, reflect.Array, reflect.Chan, reflect.Func, reflect.Interface,
			reflect.Map, reflect.Ptr, reflect.Slice, reflect.UnsafePointer, reflect.Complex64, reflect.Complex128:
			return nil, fmt.Errorf("invalid embedded struct for %s's field %s, should be struct, but got %v", schema.Name, fieldStruct.Name, fieldType)
		}
	}

	return nil, nil
}

// https://github.com/go-gorm/gorm/blob/v1.23.8/schema/schema.go#L305
func getOrParse(dest interface{}) (*Schema, error) {
	modelType := reflect.ValueOf(dest).Type()
	for modelType.Kind() == reflect.Slice || modelType.Kind() == reflect.Array || modelType.Kind() == reflect.Ptr {
		modelType = modelType.Elem()
	}

	if modelType.Kind() != reflect.Struct {
		if modelType.PkgPath() == "" {
			return nil, fmt.Errorf("%w: %+v", ErrUnsupportedDataType, dest)
		}
		return nil, fmt.Errorf("%w: %s.%s", ErrUnsupportedDataType, modelType.PkgPath(), modelType.Name())
	}

	return GetSchema(dest)
}
