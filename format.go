package gotools

import (
	"bytes"
	"encoding/json"

	"github.com/joker-circus/gotools/internal"
)

func Json(data interface{}) string {
	if v, ok := data.(string); ok {
		return v
	}

	if v, ok := data.([]byte); ok {
		return internal.B2s(v)
	}

	return JsonStruct(data)
}

func JsonIndent(data interface{}) string {
	if v, ok := data.(string); ok {
		return JsonIndentBytes([]byte(v))
	}

	if v, ok := data.([]byte); ok {
		return JsonIndentBytes(v)
	}

	return JsonIndentStruct(data)
}

func JsonStruct(data interface{}) string {
	out, _ := json.Marshal(data)
	return string(out)
}

func JsonIndentStruct(data interface{}) string {
	out, _ := json.MarshalIndent(data, "", " ")
	return string(out)
}

func JsonIndentBytes(data []byte) string {
	var b bytes.Buffer
	_ = json.Indent(&b, data, "", "    ")
	return b.String()
}

// 在JSON引号字符串中不转义有问题的HTML字符。
func JsonStructDisableEscapeHTML(data interface{}) string {
	bf := bytes.NewBuffer([]byte{})
	jsonEncoder := json.NewEncoder(bf)
	jsonEncoder.SetEscapeHTML(false)
	_ = jsonEncoder.Encode(data)
	return bf.String()
}
