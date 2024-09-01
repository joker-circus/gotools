package gotools

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type _Example struct {
	Name      string                 `json:"name"`
	Age       uint64                 `json:"age"`
	Num       *uint                  `json:"number"`
	de        string                 `json:"de"`
	Labels    map[string]interface{} `json:"labels"`
	CreatedAt time.Time              `json:"created_at"`

	d _Example2 `json:"d"`

	*_Example2

	*Example2
}

type _Example2 struct {
	Values string `json:"values"`
	Value2 string `json:"value_2"`
}
type Example2 struct {
	Values string `json:"upper_values"`
	Value2 string `json:"upper_value_2"`
}

var example _Example = _Example{
	d: _Example2{
		"d1",
		"d2",
	},
	_Example2: &_Example2{
		"ptr_private_1",
		"ptr_private_2",
	},
	Example2: &Example2{
		"ptr_1",
		"ptr_2",
	},
}

func TestStructTagAllFields(t *testing.T) {
	result := []string{"name", "age", "number", "de", "labels", "created_at", "d", "values", "value_2", "upper_values", "upper_value_2"}
	fields, _ := StructTagAllFields(example, "json")
	assert.Equal(t, result, fields, "the should be equal")
}

func TestStructTagExportedFields(t *testing.T) {
	result := []string{"name", "age", "number", "labels", "created_at", "upper_values", "upper_value_2"}
	fields, _ := StructTagExportedFields(example, "json")
	assert.Equal(t, result, fields, "the should be equal")
}

func TestStructTagExportedFieldValues(t *testing.T) {
	result := []interface{}{example.Name, example.Age, uint(0), example.Labels, example.CreatedAt, example.Example2.Values, example.Example2.Value2}
	values, _ := StructTagExportedFieldValues(example, "json")
	assert.Equal(t, result, values, "the should be equal")
}
