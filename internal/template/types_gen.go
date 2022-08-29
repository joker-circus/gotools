package template

import (
	"bytes"
	"fmt"
	"go/format"
	"io/ioutil"
	"os"
	"path"
	"runtime"
	"strings"
)

type GENERIC_NAME string

const (
	Float32 GENERIC_NAME = "Float32"
	Float64 GENERIC_NAME = "Float64"
	Int32   GENERIC_NAME = "Int32"
	Int16   GENERIC_NAME = "Int16"
	Int     GENERIC_NAME = "Int"
	Uint64  GENERIC_NAME = "Uint64"
	Uint32  GENERIC_NAME = "Uint32"
	Uint16  GENERIC_NAME = "Uint16"
	Uint    GENERIC_NAME = "Uint"
	String  GENERIC_NAME = "String"
)

var allTypes = []GENERIC_NAME{
	Float32,
	Float64,
	Int32,
	Int16,
	Int,
	Uint64,
	Uint32,
	Uint16,
	Uint,
	String,
}

// templateName，模板文件名。
// 模板文件名需要在执行调用该方法的包下，仅传递文件名即可，不需要完整路径。
// 模板文件需要用 PACKAGE_NAME、GENERIC_NAME、GENERIC_TYPE 分别代替报名、类型名、Go 类型。
//
// fileNameTemplate，生成的各类型文件名，例如："%sset.go"
// 值中%s 表示生成的 go 类型， 若无则会内容覆盖到同一文件中。
//
// types 表示对应的 Go 类型别名。不传值时，默认全部继承
func TypesGen(templateName, fileNameTemplate string, types ...GENERIC_NAME) {
	if len(types) == 0 {
		types = allTypes
	}

	_, file, _, ok := runtime.Caller(1)
	if !ok {
		panic("runtime.Caller(1) fail")
	}

	workDir := path.Dir(file)
	packageName := path.Base(workDir)

	// read template data
	f, err := os.Open(path.Join(workDir, templateName))
	if err != nil {
		panic(err)
	}
	fileData, err := ioutil.ReadAll(f)
	if err != nil {
		panic(err)
	}

	for _, v := range types {
		genericName := string(v)
		var w bytes.Buffer
		genericType := strings.ToLower(genericName)
		data := string(fileData)

		data = strings.ReplaceAll(data, "PACKAGE_NAME", packageName)
		data = strings.ReplaceAll(data, "GENERIC_NAME", genericName)
		data = strings.ReplaceAll(data, "GENERIC_TYPE", genericType)
		w.WriteString(data)

		out, err := format.Source(w.Bytes())
		if err != nil {
			panic(err)
		}

		filename := path.Join(workDir, fmt.Sprintf(fileNameTemplate, genericType))
		if err := ioutil.WriteFile(filename, out, 0660); err != nil {
			panic(err)
		}
	}
}
