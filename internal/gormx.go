package internal

import "strings"

// 获取 Gorm Tag 中 Column 值。
func GetGormTagColumnName(str string) (v string, ok bool) {
	tags := ParseTagSetting(str, ";")
	v, ok = tags["COLUMN"]
	if ok {
		return
	}
	if len(tags) == 1 {
		return str, len(str) != 0
	}
	return
}

// https://github.com/go-gorm/gorm/blob/v1.23.8/schema/utils.go#L16
func ParseTagSetting(str string, sep string) map[string]string {
	settings := map[string]string{}
	names := strings.Split(str, sep)

	for i := 0; i < len(names); i++ {
		j := i
		if len(names[j]) > 0 {
			for {
				if names[j][len(names[j])-1] == '\\' {
					i++
					names[j] = names[j][0:len(names[j])-1] + sep + names[i]
					names[i] = ""
				} else {
					break
				}
			}
		}

		values := strings.Split(names[j], ":")
		k := strings.TrimSpace(strings.ToUpper(values[0]))

		if len(values) >= 2 {
			settings[k] = strings.Join(values[1:], ":")
		} else if k != "" {
			settings[k] = k
		}
	}

	return settings
}
