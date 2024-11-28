package datatypes

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
)

type SerializerInterface interface {
	sql.Scanner
	driver.Valuer
	json.Marshaler
	json.Unmarshaler
}

var _ SerializerInterface = &DBJsonField[string]{}

type DBJsonField[T any] struct {
	field T
}

func NewDBJsonField[T any](data T) DBJsonField[T] {
	return DBJsonField[T]{
		field: data,
	}
}

func (f DBJsonField[T]) Data() T {
	return f.field
}

func (f DBJsonField[T]) Value() (driver.Value, error) {
	b, err := f.MarshalJSON()
	if err != nil {
		return string(b), err
	}

	if len(b) == 4 && string(b) == "null" {
		return "", nil
	}
	return string(b), nil
}

func (f *DBJsonField[T]) Scan(v interface{}) error {
	var val DBJsonField[T]
	if v == nil {
		*f = val
		return nil
	}

	bytes, ok := v.([]byte)
	if !ok {
		return fmt.Errorf("can not scan value %v to %T", v, val.field)
	}

	if len(bytes) == 0 {
		*f = val
		return nil
	}

	return f.UnmarshalJSON(bytes)
}

func (f DBJsonField[T]) MarshalJSON() ([]byte, error) {
	return json.Marshal(f.field)
}

func (f *DBJsonField[T]) UnmarshalJSON(b []byte) error {
	return json.Unmarshal(b, &f.field)
}

// 存储 DB 时使用逗号拼接
type DBStringSlice []string

func (f DBStringSlice) Value() (driver.Value, error) {
	return f.ToString(), nil
}

func (f DBStringSlice) ToString() string {
	if len(f) == 0 {
		return ""
	}
	return strings.Join(f, ",")
}

func (f *DBStringSlice) Scan(v interface{}) error {
	if v == nil {
		*f = []string{}
		return nil
	}

	bytes, ok := v.([]byte)
	if !ok {
		return fmt.Errorf("can not scan value %v to %T", v, *f)
	}

	if len(bytes) == 0 {
		*f = []string{}
		return nil
	}

	return f.UnmarshalJSON(bytes)
}

func (f DBStringSlice) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf(`"%s"`, strings.Join(f, ","))), nil
}

func (f *DBStringSlice) UnmarshalJSON(b []byte) error {
	*f = strings.Split(strings.Trim(string(b), `"`), ",")
	return nil
}

// Store in DB using comma concatenation
type DBFloat64Slice []float64

func (f DBFloat64Slice) Value() (driver.Value, error) {
	return f.ToString(), nil
}

func (f DBFloat64Slice) ToString() string {
	if len(f) == 0 {
		return ""
	}
	stringSlice := make([]string, 0, len(f))
	for _, v := range f {
		stringSlice = append(stringSlice, strconv.FormatFloat(v, 'f', -1, 64))
	}
	return strings.Join(stringSlice, ",")
}

func (f *DBFloat64Slice) Scan(v interface{}) error {
	if v == nil {
		*f = []float64{}
		return nil
	}

	bytes, ok := v.([]byte)
	if !ok {
		return fmt.Errorf("can not scan value %v to %T", v, *f)
	}

	if len(bytes) == 0 {
		*f = []float64{}
		return nil
	}

	return f.UnmarshalJSON(bytes)
}

func (f DBFloat64Slice) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf(`"%s"`, f.ToString())), nil
}

func (f *DBFloat64Slice) UnmarshalJSON(b []byte) error {
	stringSlice := strings.Split(strings.Trim(string(b), `"`), ",")
	float64Slice := make([]float64, 0, len(stringSlice))
	for _, v := range stringSlice {
		vv, _ := strconv.ParseFloat(v, 64)
		float64Slice = append(float64Slice, vv)
	}
	*f = float64Slice
	return nil
}
