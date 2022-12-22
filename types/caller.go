package types

import (
	"path"
	"reflect"
	"runtime"
)

// 获取方法名
func NameOfFunction(f interface{}) string {
	return runtime.FuncForPC(reflect.ValueOf(f).Pointer()).Name()
}

// 获取调用函数的 Caller 信息。
// 即 A 调用 B，B 通过 CallerInfo() 可以获取 A 的信息。
//
// 返回：A 方法名、A 的文件名、A 的包名、A 的工作路径、A 调用 B 的所在行号、是否获取到 caller 信息。
func CallerInfo() (funcName, fileName, packageName, workDir string, line int, ok bool) {
	return getCallerInfo(3)
}

// 获取自身的函数的 Caller 信息。
// 即 A 调用 SelfCallerInfo() 可以获取 A 的信息。
//
// 返回：A 方法名、A 的文件名、A 的包名、A 的工作路径、A 调用 B 的所在行号、是否获取到 caller 信息。
func SelfCallerInfo() (funcName, fileName, packageName, workDir string, line int, ok bool) {
	return getCallerInfo(2)
}

// 获取调用者的 Caller 信息。
//
// 返回：方法名、文件名、包名、工作路径、执行函数的所在行号、是否获取到 caller 信息。
func getCallerInfo(skip int) (funcName, fileName, packageName, workDir string, line int, ok bool) {
	pc, file, line, ok := runtime.Caller(skip)
	if !ok {
		return
	}

	fileName = path.Base(file)
	workDir = path.Dir(file)
	packageName = path.Base(workDir)
	funcName = runtime.FuncForPC(pc).Name()
	return
}
