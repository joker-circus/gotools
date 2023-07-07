package gotools

import (
	"strconv"
	"strings"
	"unicode"

	"github.com/joker-circus/gotools/internal"
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

// 首字母大写
func Ucfirst(str string) string {
	for i, v := range str {
		return string(unicode.ToUpper(v)) + str[i+1:]
	}
	return ""
}

// 首字母小写
func Lcfirst(str string) string {
	for i, v := range str {
		return string(unicode.ToLower(v)) + str[i+1:]
	}
	return ""
}

// 驼峰写法。下划线写法转为驼峰写法。
//
// 例如：xx_yy to XxYx  xx_y_y to XxYY。
//
// 等同于：
//
//	s = strings.Replace(s, "_", " ", -1)
//	s = strings.Title(s)
//	return strings.Replace(s, " ", "", -1)
func CamelCase(s string) string {
	data := make([]byte, 0, len(s))
	j := false
	k := false
	num := len(s) - 1
	for i := 0; i <= num; i++ {
		d := s[i]
		if k == false && d >= 'A' && d <= 'Z' {
			k = true
		}
		if d >= 'a' && d <= 'z' && (j || k == false) {
			d = d - 32
			j = false
			k = true
		}
		if k && d == '_' && num > i && s[i+1] >= 'a' && s[i+1] <= 'z' {
			j = true
			continue
		}
		data = append(data, d)
	}
	return string(data[:])
}

// 蛇形/下划线写法。驼峰式写法转为下划线写法。
//
// 例如：XxYy to xx_yy , XxYY to xx_y_y。
func SnakeCase(s string) string {
	data := make([]byte, 0, len(s)*2)
	j := false
	num := len(s)
	for i := 0; i < num; i++ {
		d := s[i]
		// or通过ASCII码进行大小写的转化
		// 65-90（A-Z），97-122（a-z）
		//判断如果字母为大写的A-Z就在前面拼接一个_
		if i > 0 && d >= 'A' && d <= 'Z' && j {
			data = append(data, '_')
		}
		if d != '_' {
			j = true
		}
		data = append(data, d)
	}
	//ToLower把大写字母统一转小写
	return strings.ToLower(string(data[:]))
}

// Domain get the domain of given URL
func Domain(url string) string {
	domainPattern := `([a-z0-9][-a-z0-9]{0,62})\.` +
		`(com\.cn|com\.hk|` +
		`cn|com|net|edu|gov|biz|org|info|pro|name|xxx|xyz|be|` +
		`me|top|cc|tv|tt)`
	domain := internal.MatchOneOf(url, domainPattern)
	if domain != nil {
		return domain[1]
	}
	return ""
}

// LimitLength Handle overly long strings
func LimitLength(s string, length int) string {
	// 0 means unlimited
	if length == 0 {
		return s
	}

	const ELLIPSES = "..."
	str := []rune(s)
	if len(str) > length {
		return string(str[:length-len(ELLIPSES)]) + ELLIPSES
	}
	return s
}

// Reverse Reverse a string
func Reverse(s string) string {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}
