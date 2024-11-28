package timeutil

import (
	"fmt"
	"log"
	"sync"
	"time"
)

type MeasureLogFunc func(format string, args ...interface{})

var defaultLog MeasureLogFunc = func(format string, args ...interface{}) {
	log.Println(fmt.Sprintf(format, args...))
}

var enableMeasure bool

// 设置测量函数日志。
// 默认：fmt.Println(fmt.Sprintf(format, args...))。
func SetMeasureLog(f MeasureLogFunc) {
	defaultLog = f
}

// 是否启用测量函数。
// 默认 false。
func EnableMeasureTime(enable bool) {
	enableMeasure = enable
}

// 测量动作用时时长。
// 使用方法：defer MeasureTime(actionName)()。
// 末尾必须带上"()"！！！。
func MeasureTime(actionName string) func() {
	if !enableMeasure {
		return func() {}
	}

	start := time.Now()
	defaultLog("%s action is beginning...", actionName)
	return func() {
		defaultLog("Time taken by %s action is %v", actionName, time.Since(start))
	}
}

// 多次测量动作用时时长。
// 参数 end 表示是否结束，结束后会计算总的用时时长，后续再调用无效。
//
// 使用方法：
//
//	measure := GetMeasurer(父动作名)
//	……
//	measure(子动作名，false)
//	……
//	measure(最后的子动作名，true)
func GetMeasurer(actionName string) func(subActionName string, end bool) {
	var realEnd bool
	start := time.Now()
	t1 := start
	defaultLog("%s.action.is.running", actionName)
	return func(subActionName string, end bool) {
		if realEnd {
			return
		}
		t2 := time.Now()
		if subActionName != "" || !end {
			defaultLog("Time.taken.by.%s.%s.action.is: %v", actionName, subActionName, t2.Sub(t1))
		}
		if end {
			defaultLog("Time.taken.by.%s.action.is: %v", actionName, t2.Sub(start))
		}
		t1 = t2
		realEnd = end
	}
}

type Measurer struct {
	actionName          string
	start, preTime, end time.Time
	finish              bool
	defaultLog          MeasureLogFunc
	m                   sync.RWMutex
}

func NewMeasurer(actionName string) *Measurer {
	start := time.Now()
	defaultLog("%s.action.is.running", actionName)
	return &Measurer{
		actionName: actionName,
		start:      start,
		preTime:    start,
		finish:     false,
		defaultLog: defaultLog,
	}
}

func (m *Measurer) SetLog(f MeasureLogFunc) {
	m.m.Lock()
	m.defaultLog = f
	m.m.Unlock()
}

func (m *Measurer) Measure(subActionName string) {
	if m.finish {
		return
	}
	end := time.Now()
	if subActionName == "" {
		subActionName = "null"
	}
	m.defaultLog("Time.taken.by.%s.%s.action.is: %v", m.actionName, subActionName, end.Sub(m.preTime))

	m.m.Lock()
	m.preTime = end
	m.m.Unlock()
}

func (m *Measurer) AsyncMeasure(subActionName string, subStart time.Time) {
	if m.finish {
		return
	}
	end := time.Now()
	if subActionName == "" {
		subActionName = "null"
	}
	m.defaultLog("Time.taken.by.%s.%s.action.is: %v", m.actionName, subActionName, end.Sub(subStart))

	m.m.Lock()
	m.preTime = end
	m.m.Unlock()
}

func (m *Measurer) EndMeasure() {
	if m.finish {
		return
	}
	end := time.Now()
	m.defaultLog("Time.taken.by.%s.action.is: %v", m.actionName, end.Sub(m.start))

	m.m.Lock()
	m.finish = true
	m.preTime = end
	m.end = end
	m.m.Unlock()
}
