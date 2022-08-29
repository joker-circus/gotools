package datatypes

import (
   	"fmt"
	"time"
	"database/sql/driver"
)

// Time is alias type for time.Time
type Time time.Time

const (
	timeFormat = "2006-01-02 15:04:05"
	dateFormat = "2006-01-02"
	zone       = "Asia/Shanghai"
)

var (
	zeroTime, _ = time.Parse(time.RFC3339, "1970-01-01T08:00:00+08:00")
)

// UnmarshalJSON implements json unmarshal interface.
func (t *Time) UnmarshalJSON(data []byte) (err error) {
	now, err := time.ParseInLocation(`"`+timeFormat+`"`, string(data), time.Local)
  if err != nil {
		now, err = time.ParseInLocation(`"`+time.RFC3339+`"`, string(data), time.Local)
	}
	*t = Time(now)
	return
}

// MarshalJSON implements json marshal interface.
func (t Time) MarshalJSON() ([]byte, error) {
	b := make([]byte, 0, len(timeFormat)+2)
	b = append(b, '"')
	b = time.Time(t).AppendFormat(b, timeFormat)
	b = append(b, '"')
	return b, nil
}

func (t Time) String() string {
	return time.Time(t).Format(timeFormat)
}

func (t Time) DisplayTimeString() string {
	if t.IsZero() {
		return "-"
	}

	return time.Time(t).Format(timeFormat)
}

func (t Time) DisplayDateString() string {
	if t.IsZero() {
		return "-"
	}

	return time.Time(t).Format(dateFormat)
}

func (t Time) local() time.Time {
	loc, _ := time.LoadLocation(zone)
	return time.Time(t).In(loc)
}

// IsInitialized 检查是否已经被正确赋值过
func (t Time) IsInitialized() bool {
	return !time.Time(t).IsZero()
}

func Now() Time {
	return Time(time.Now())
}

func (t Time) Value() (driver.Value, error) {
	if !t.IsInitialized() {
		return zeroTime, nil
	}
	return time.Time(t), nil
}

// Scan 和 Value 必须使用 time.Time 作为输入输出以兼容对 time.Time 类型的默认操作
func (t *Time) Scan(v interface{}) error {
	if v == nil {
		*t = Time(zeroTime)
		return nil
	}
	value, ok := v.(time.Time)
	if !ok {
		return fmt.Errorf("can not scan value %v to timeconv", v)
	}
	*t = Time(value)
	return nil
}

func (t Time) Before(t1 Time) bool {
	tt := time.Time(t)
	tt1 := time.Time(t1)

	return tt.Before(tt1)
}

func (t Time) After(t1 Time) bool {
	tt := time.Time(t)
	tt1 := time.Time(t1)

	return tt.After(tt1)
}

func ZeroTime() Time {
	return Time(zeroTime)
}

func (t Time) IsZero() bool {
	return time.Time(t).Equal(zeroTime)
}

func (t Time) Equal(t1 Time) bool {
	return time.Time(t).Equal(time.Time(t1))
}

func (t Time) Unix() int64 {
	tt := time.Time(t)
	return tt.Unix()
}

func (t Time) AddDate(years int, months int, days int) Time {
	tt := time.Time(t).AddDate(years, months, days)
	return Time(tt)
}

func (t Time) Format(layout string) string {
	return time.Time(t).Format(layout)
}
