package template

import (
	"bytes"
	"fmt"
	"strings"
	"text/template"
)

func Render(templateContent string, varsMap interface{}) (string, error) {
	// 渲染
	tmpl, err := template.New("default").Option("missingkey=error").Parse(templateContent)
	if err != nil {
		return "", err
	}
	buf := new(bytes.Buffer)
	err = tmpl.Execute(buf, varsMap)
	if err != nil {
		return "", RenderErrFormat(err)
	}

	// 替换文本中 {{ if }} 等语句带来的多行
	return strings.ReplaceAll(buf.String(), "\n\n", "\n"), nil
}

// simplify error
func RenderErrFormat(err error) error {
	if strings.Contains(err.Error(), "map has no entry") {
		tmpSlice := strings.Split(err.Error(), " ")
		missingKey := tmpSlice[len(tmpSlice)-1]
		missingKey = "[" + missingKey[1:len(missingKey)-1] + "]"
		return fmt.Errorf("%s was not defined", missingKey)
	}
	return err
}