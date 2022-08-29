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
