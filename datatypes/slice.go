package datatypes

import (
	"encoding/json"
	"fmt"

	"database/sql/driver"
)

type Slice[T any] []T

func (s Slice[T]) Value() (driver.Value, error) {
	bytes, err := json.Marshal(s)
	if err != nil {
		return nil, err
	}

	return bytes, nil
}

func (s *Slice[T]) Scan(v interface{}) error {
	if v == nil {
		return nil
	}

	bytes, ok := v.([]byte)
	if !ok {
		return fmt.Errorf("can not scan value %v to Slice", v)
	}

	if len(bytes) == 0 {
		return nil
	}

	var value []T
	err := json.Unmarshal(bytes, &value)
	if err != nil {
		return err
	}

	*s = value
	return nil
}
