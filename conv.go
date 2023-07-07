package gotools

import (
	"fmt"
	"reflect"
	"strconv"
	"time"
)

func ToString(v interface{}) string {
	switch value := v.(type) {
	case string:
		return value
	case bool:
		return strconv.FormatBool(value)
	case int:
		return strconv.Itoa(value)
	case int8:
		return strconv.FormatInt(int64(value), 10)
	case int16:
		return strconv.FormatInt(int64(value), 10)
	case int32:
		return strconv.FormatInt(int64(value), 10)
	case int64:
		return strconv.FormatInt(value, 10)
	case uint:
		return strconv.FormatUint(uint64(value), 10)
	case uint8:
		return strconv.FormatUint(uint64(value), 10)
	case uint16:
		return strconv.FormatUint(uint64(value), 10)
	case uint32:
		return strconv.FormatUint(uint64(value), 10)
	case uint64:
		return strconv.FormatUint(value, 10)
	case float32:
		return strconv.FormatFloat(float64(value), 'f', -1, 32)
	case float64:
		return strconv.FormatFloat(value, 'f', -1, 64)
	case fmt.Stringer:
		return value.String()
	default:
		return fmt.Sprintf("%v", v)
	}
}

func ToFloat64(v interface{}) (float64, error) {
	switch value := v.(type) {
	case float32:
		return float64(value), nil
	case float64:
		return value, nil
	case int:
		return float64(value), nil
	case int8:
		return float64(value), nil
	case int16:
		return float64(value), nil
	case int32:
		return float64(value), nil
	case int64:
		return float64(value), nil
	case uint:
		return float64(value), nil
	case uint8:
		return float64(value), nil
	case uint16:
		return float64(value), nil
	case uint32:
		return float64(value), nil
	case uint64:
		return float64(value), nil
	case string:
		f, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return 0, fmt.Errorf("cannot convert string '%s' to float64", value)
		}
		return f, nil
	default:
		return 0, fmt.Errorf("cannot convert type %v to float64", reflect.TypeOf(v))
	}
}

var defaultLayouts = []string{
	"2006-01-02 15:04:05",
	"2006-01-02",
	"2006/01/02 15:04:05",
	"2006/01/02",
	"01/02/2006 15:04:05",
	"01/02/2006",
	time.Layout,
	time.ANSIC,
	time.UnixDate,
	time.RubyDate,
	time.RFC822,
	time.RFC822Z,
	time.RFC850,
	time.RFC1123,
	time.RFC1123Z,
	time.RFC3339,
	time.RFC3339Nano,
}

func Time(v interface{}, layouts ...string) (time.Time, bool) {
	switch value := v.(type) {
	case time.Time:
		return value, true
	case string:
		for _, layout := range append(defaultLayouts, layouts...) {
			t, err := time.Parse(layout, value)
			if err == nil {
				return t, true
			}
		}
	}
	return time.Now(), false
}
