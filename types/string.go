package types

import (
	"strconv"
	"strings"
)

func LowerSlice(values ...string) []string {
	n := make([]string, len(values))
	for i, v := range values {
		n[i] = strings.ToLower(v)
	}
	return n
}

func InSlice(s []string, val string) bool {
	for _, v := range s {
		if v == val {
			return true
		}
	}
	return false
}

// 判断字符串是不是 true
func IsTrueString(str string) bool {
	return strings.ToLower(str) == "true"
}

func HumanUnit(l uint64) string {
	var (
		suffix string
		size   float64
	)

	if l >= (1 << 30) {
		suffix = "G"
		size = float64(l) / (1 << 30)
	} else if l >= (1 << 20) {
		suffix = "M"
		size = float64(l) / (1 << 20)
	} else if l >= (1 << 10) {
		suffix = "K"
		size = float64(l) / (1 << 10)
	} else {
		size = float64(l)
		suffix = "B"
	}

	return strconv.FormatFloat(size, 'f', 3, 64) + suffix
}

