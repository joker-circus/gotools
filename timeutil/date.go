package timeutil

import "time"

const (
	OneDayUnixSecond      = 86400
	OneDayUnixMilliSecond = 86400000
)

// 一天起始时间
func GetTodayUnixSecond() uint64 {
	return GetUnixSecond(time.Now())
}

// 一天起始时间
func GetTodayUnixMilliSecond() uint64 {
	return GetUnixMilliSecond(time.Now())
}

func GetUnixSecond(timeObj time.Time) uint64 {
	formatStr := "2006-01-02"
	return GetUnixSecondFromStr(timeObj.Format(formatStr))
}

// date format "2006-01-02"
func GetUnixSecondFromStr(date string) uint64 {
	formatStr := "2006-01-02"
	if len(date) > len(formatStr) {
		date = date[:len(formatStr)]
	}
	today, _ := time.ParseInLocation(formatStr, date, time.Local)
	return uint64(today.Unix())
}

func GetUnixMilliSecond(timeObj time.Time) uint64 {
	return GetUnixSecond(timeObj) * 1000
}

// fmtDate 格式化 date 日期字符串，返回只精确到天的 Time 值。
func fmtDate(date string) time.Time {
	layout := "2006-01-02"
	if len(date) < len(layout) {
		return time.Time{}
	}
	if len(date) > len(layout) {
		date = date[:len(layout)]
	}
	t, _ := time.Parse("2006-01-02", date)
	return t
}

// RangeDate 遍历 startDate 到 endDate 之间的日期，
// startDate 和 endDate 要求 “2006-01-02” 格式，
// f 接收遍历的日期值进行操作，返回 true 表示继续遍历，否则停止。
func RangeDate(startDate, endDate string, f func(date string) bool) {
	start, end := fmtDate(startDate), fmtDate(endDate)
	for !start.After(end) {
		if !f(start.Format("2006-01-02")) {
			return
		}
		start = start.AddDate(0, 0, 1)
	}
}
