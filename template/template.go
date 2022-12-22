package template

import (
	"bytes"
	"fmt"
	"regexp"
	"strings"
	"text/template"
)

// RenderFuncMap 是 Render 系列方法默认的 Template.FuncMap
var RenderFuncMap template.FuncMap

// StringsFuncMap 提供 strings 的 Contains、Replace、HasSuffix、HasPrefix、Title、
// ToTitle、ToLower、ToUpper、Split、Trim、TrimSpace、Join方法
var StringsFuncMap = template.FuncMap{
	"Contains":  strings.Contains,
	"Replace":   strings.ReplaceAll,
	"HasSuffix": strings.HasSuffix,
	"HasPrefix": strings.HasPrefix,
	"Title":     strings.Title,
	"ToTitle":   strings.ToTitle,
	"ToLower":   strings.ToLower,
	"ToUpper":   strings.ToUpper,
	"Split":     strings.Split,
	"Trim":      strings.Trim,
	"TrimSpace": strings.TrimSpace,
	"Join":      strings.Join,
}

func init() {
	RenderFuncMap = StringsFuncMap
}

// Render 默认替换文本中 {{ if }} 等语句带来的多行
func Render(templateContent string, varsMap interface{}) (string, error) {
	return RenderWithOptions(templateContent, varsMap, ReplaceMultipleLine)
}

func RenderWithOptions(content string, varsMap interface{}, options ...RenderOption) (string, error) {
	result, err := render(content, varsMap, RenderFuncMap)
	if err != nil {
		return "", err
	}

	for _, option := range options {
		result = option(result)
	}

	return result, nil
}

func render(content string, varsMap interface{}, funcMap template.FuncMap) (string, error) {
	// 渲染
	tmpl, err := template.New("default").Funcs(funcMap).Option("missingkey=error").Parse(content)
	if err != nil {
		return "", err
	}
	buf := new(bytes.Buffer)
	err = tmpl.Execute(buf, varsMap)
	if err != nil {
		return "", renderErrFormat(err)
	}

	return buf.String(), nil
}

// simplify error
func renderErrFormat(err error) error {
	if strings.Contains(err.Error(), "map has no entry") {
		tmpSlice := strings.Split(err.Error(), " ")
		missingKey := tmpSlice[len(tmpSlice)-1]
		missingKey = "[" + missingKey[1:len(missingKey)-1] + "]"
		return fmt.Errorf("%s was not defined", missingKey)
	}
	return err
}

// RenderOption 用于对渲染后的文本信息进行额外处理
type RenderOption func(string) string

// 等同于 “\\n+” 正则，二次转义
var multipleLineReg = regexp.MustCompile(`\n+`)

// 替换多余的空行
func ReplaceMultipleLine(s string) string {
	return multipleLineReg.ReplaceAllString(s, "\n")
}

// 替换 & -> &amp;
func ReplaceAmp(s string) string {
	return strings.ReplaceAll(s, "&", "&amp;")
}
