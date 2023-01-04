package types

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

func TestStructTagAllFields(t *testing.T) {
	var example _Example
	result := []string{"name", "age", "number", "de", "labels", "created_at", "d", "values", "value_2", "upper_values", "upper_value_2"}
	tags, _ := StructTagAllFields(example, "json")
	assert.Equal(t, result, tags, "the should be equal")
}

func TestStructTagExportedFields(t *testing.T) {
	var example _Example
	result := []string{"name", "age", "number", "labels", "created_at", "upper_values", "upper_value_2"}
	tags, _ := StructTagExportedFields(example, "json")
	assert.Equal(t, result, tags, "the should be equal")
}
