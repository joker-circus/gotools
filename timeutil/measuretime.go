package timeutil

import (
	"fmt"
	"time"
)

type MeasureLogFunc func(format string, args ...interface{})

var defaultLog MeasureLogFunc = func(format string, args ...interface{}) {
	fmt.Printf(format, args...)
}

var enableMeasure bool

//	设置测量函数日志
// 	默认：fmt.Printf(format, args...)
func SetMeasureLog(f MeasureLogFunc) {
	defaultLog = f
}

// 	是否启用测量函数
// 	默认 false
func EnableMeasureTime(enable bool) {
	enableMeasure = enable
}

// 	测量动作用时时长
// 	使用方法：defer MeasureTime(actionName)()
// 	末尾必须带上"()"！！！
func MeasureTime(actionName string) func() {
	if !enableMeasure {
		return func() {}
	}

	start := time.Now()
	return func() {
		defaultLog("Time taken by %s action is %v \n", actionName, time.Since(start))
	}
}
