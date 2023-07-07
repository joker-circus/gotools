package gotools

import "bytes"

// 不区分大小写比较
func CaseInsensitiveCompare(a, b []byte) bool {
	if len(a) != len(b) {
		return false
	}
	for i := 0; i < len(a); i++ {
		if a[i]|0x20 != b[i]|0x20 {
			return false
		}
	}
	return true
}

// 拆分行
func NextLine(b []byte) ([]byte, []byte, bool) {
	nNext := bytes.IndexByte(b, '\n')
	if nNext < 0 {
		return nil, nil, false
	}
	n := nNext
	if n > 0 && b[n-1] == '\r' {
		n--
	}
	return b[:n], b[nNext+1:], true
}
