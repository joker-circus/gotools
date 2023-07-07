package timeutil

import (
	"errors"
	"fmt"
	"regexp"
	"time"
	"unsafe"
)

const (
	DateTimeFormat = "2006-01-02 15:04:05"
	DateFormat     = "2006-01-02"
	MonthFormat    = "2006-01"
)

var nullTime time.Time

// CompareDay is 比较天级的时间，如果first >= second 返回true,不然返回false
func CompareDay(first, second time.Time) bool {
	fYear, fMonth, fDay := first.Date()
	sYear, sMonth, sDay := second.Date()
	if fYear < sYear {
		return false
	}
	if fMonth < sMonth {
		return false
	}
	if fDay < sDay {
		return false
	}
	return true
}

// ToDay is 转换时间为  年月日，转换后格式为  xxxx-xx-xx 00:00:00 对应的时间戳
func ToDay(origin time.Time) int64 {
	newTime := origin.Format(DateFormat)
	loc, _ := time.LoadLocation("Local")
	dt, _ := time.ParseInLocation(DateFormat, newTime, loc)
	return dt.Unix()
}

// ToDayString, 将时间字符串转换为time.Time
func ToDayString(origin string) (time.Time, error) {
	loc, _ := time.LoadLocation("Local")
	dt, err := time.ParseInLocation(DateFormat, origin, loc)
	if err != nil {
		return nullTime, err
	}
	return dt, nil
}

// ToDayString, 将时间字符串转换为time.Time
func ToDateTime(origin string) (time.Time, error) {
	loc, _ := time.LoadLocation("Local")
	dt, err := time.ParseInLocation(DateTimeFormat, origin, loc)
	if err != nil {
		return nullTime, err
	}
	return dt, nil
}

// ToMonthTime, 将时间字符串转换为time.Time,传入的时间必须是月级别，例如2020-01
func ToMonthTime(origin string) (time.Time, error) {
	loc, _ := time.LoadLocation("Local")
	dt, err := time.ParseInLocation(MonthFormat, origin, loc)
	if err != nil {
		return nullTime, err
	}
	return dt, nil
}

// 算出两个时间差值，返回天数
func GetTimeSub(at, bt time.Time) (int64, error) {
	if at.After(bt) {
		return 0, errors.New("startTime should be earlier than EndTime")
	}

	startTime := at.Unix()
	endTime := bt.Unix()
	// 求相差天数
	date := (endTime - startTime) / 86400
	return date, nil

}

// ParseDateRange ...
// 传入的值必须是 2020-08-10 格式的日期，第一个参数是起始日期，第二个参数是截止日期，s 比 e 前
// 解析出来的是 s 和 e 之间的日期列表的time 格式。
func ParseDateRange(s, e string) ([]time.Time, error) {
	// 校验格式
	r, err := regexp.Compile(`([\d][\d][\d][\d])-([\d][\d])-([\d][\d])$`)
	if err != nil {
		return nil, fmt.Errorf("regexp compile err, err: %s", err.Error())
	}
	if !r.MatchString(s) {
		return nil, errors.New("stime 格式不正确")
	}
	if !r.MatchString(e) {
		return nil, errors.New("etime 格式不正确")
	}

	loc, err := time.LoadLocation("Local")
	if err != nil {
		return nil, err
	}

	startTime, err := time.ParseInLocation(DateFormat, s, loc)
	if err != nil {
		return nil, err
	}

	endTime, err := time.ParseInLocation(DateFormat, e, loc)
	if err != nil {
		return nil, err
	}

	if s == e {
		return []time.Time{startTime, endTime}, nil
	}

	timeSub, err := GetTimeSub(startTime, endTime)
	if err != nil {
		return nil, err
	}

	du, err := time.ParseDuration("24h")
	if err != nil {
		return nil, err
	}

	times := make([]time.Time, *(*int)(unsafe.Pointer(&timeSub)))
	t := startTime
	for i := 0; i < *(*int)(unsafe.Pointer(&timeSub)); i++ {
		times[i] = t
		t = t.Add(du)
	}

	return times, nil
}

// TimeAdd24H ... 获取输入时间的24小时候的时间
func TimeAdd24H(st time.Time) (time.Time, error) {
	du, err := time.ParseDuration("24h")
	if err != nil {
		return st, err
	}

	return st.Add(du), nil
}

// TimeAdd24H ... 获取输入时间的24小时候的时间
func TimeDel24H(st time.Time) (time.Time, error) {
	du, err := time.ParseDuration("-24h")
	if err != nil {
		return st, err
	}

	return st.Add(du), nil
}

func GetTodayDuration(start, end time.Time, todayDate string) (int64, int64, int64, error) {
	date, err := time.ParseInLocation(DateTimeFormat, todayDate+" 00:00:00", time.Local)
	if err != nil {
		return 0, 0, 0, fmt.Errorf("error occurred when format time, error is %s, time is %s", err.Error(), todayDate)
	}
	dateToUnix := date.Unix()

	var startAt, endAt, todayDuration int64
	startAt = start.Unix()
	endAt = end.Unix()
	if dateToUnix > startAt {
		startAt = dateToUnix
	}

	if endAt <= 0 {
		endDate := date.AddDate(0, 0, 1).Unix()
		endAt = endDate
	}

	todayDuration = endAt - startAt
	return startAt, endAt, todayDuration, nil
}

// 计算某一个月的天数
func CountMonthDays(month string) (int, error) {
	monthTime, err := ToMonthTime(month)
	if err != nil {
		return 0, err
	}

	return count(monthTime.Year(), monthTime.Month()), nil
}
func count(year int, month time.Month) (days int) {
	if month != 2 {
		if month == 4 || month == 6 || month == 9 || month == 11 {
			days = 30
		} else {
			days = 31
		}
	} else {
		if year%4 == 0 && (year%100 != 0 || year%400 == 0) {
			days = 29
		} else {
			days = 28
		}
	}
	return
}

// 获取传入的时间所在月份的第一天，即某月第一天的0点。如传入time.Now(), 返回当前月份的第一天0点时间
func GetFirstDateOfMonth(d time.Time) time.Time {
	d = d.AddDate(0, 0, -d.Day()+1)
	return GetZeroTime(d)
}

// 获取传入的时间所在月份的最后一天，即某月最后一天的0点。如传入time.Now(), 返回当前月份的最后一天0点时间。
func GetLastDateOfMonth(d time.Time) time.Time {
	return GetFirstDateOfMonth(d).AddDate(0, 1, -1)
}

func GetZeroTime(d time.Time) time.Time {
	return time.Date(d.Year(), d.Month(), d.Day(), 0, 0, 0, 0, d.Location())
}

func ValidDateOrYesterday(date string) (string, error) {
	if date == "" {
		return time.Now().AddDate(0, 0, -1).Format(DateFormat), nil
	}

	if _, err := ToDayString(date); err != nil {
		return time.Now().AddDate(0, 0, -1).Format(DateFormat), err
	}

	return date, nil
}

// GetLastTimePoint 获取上一个interval的“整点”值, 截掉多余的尾数
// 比如，interval是10，t在00:10:00~00:19:59，都将返回00:10:00的时间点值
// interval 是30时，t在01:30:00~01:59:59，都将返回01:30:00的时间点值
func GetLastTimePoint(t time.Time, interval int) time.Time {
	if interval <= 0 {
		return t
	}
	tt := time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), 0, 0, t.Location())
	e := time.Duration(t.Minute() % interval)
	return tt.Add(-e * time.Minute)
}

// GetStartTimeOfDate 获取时间d当天的0点值
func GetStartTimeOfDate(d time.Time) time.Time {
	return time.Date(d.Year(), d.Month(), d.Day(), 0, 0, 0, 0, d.Location())
}

// GetEndTimeOfDate 获取时间d当天的23:59:59值
func GetEndTimeOfDate(d time.Time) time.Time {
	return time.Date(d.Year(), d.Month(), d.Day(), 23, 59, 59, 0, d.Location())
}

// RangeTime 将 t1, t2 按照最大 interval 划分。
func RangeTime(t1, t2 time.Time, interval time.Duration, f func(t1, t2 time.Time) bool) {
	var end time.Time
	for t1.Before(t2) {
		// 调用函数 f 处理打印操作
		end = t1.Add(interval)
		if end.After(t2) {
			end = t2
		}
		if !f(t1, end) {
			return
		}
		t1 = end
	}
}
