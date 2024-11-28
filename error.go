package gotools

import (
	"runtime"

	"github.com/pkg/errors"
)

// 获取堆栈信息。
// 如 A -> B -> C -> NewStackTrace，skip = 0 获取从 C 开始的堆栈，skip = 1 获取从 B 开始的堆栈，skip = 2 获取从 A 开始的堆栈，依次类推。
// 实际参考，例如：通过 fmt.Sprintf("%s\n%+v", err.Error(), NewStackTrace(1)) 可以获取到 error 所在的堆栈信息
func StackTrace(skip int) errors.StackTrace {
	const depth = 32

	var pcs [depth]uintptr
	n := runtime.Callers(2+skip, pcs[:])

	f := make(errors.StackTrace, n)
	for i := 0; i < n; i++ {
		f[i] = errors.Frame(pcs[i])
	}

	return f
}
