package format

import (
	"bytes"
	"encoding/json"
)

func Json(data interface{}) string {
	if v, ok := data.(string); ok {
		return JsonBytes([]byte(v))
	}

	if v, ok := data.([]byte); ok {
		return JsonBytes(v)
	}

	return JsonStruct(data)
}

func JsonStruct(data interface{}) string {
	out, _ := json.MarshalIndent(data, "", " ")
	return string(out)
}

func JsonBytes(data []byte) string  {
	var b bytes.Buffer
	_ = json.Indent(&b, data, "", "    ")
	return b.String()
}
