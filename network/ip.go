package network

import (
	"github.com/pkg/errors"
)

type IPV4 [4]byte

const IPLimit = 255

// Parse IPv4 address (d.d.d.d).
func parseIPV4(s string) (p [4]byte, ok bool) {
	for i := 0; i < 4; i++ {
		if len(s) == 0 {
			// Missing octets.
			return p, false
		}
		if i > 0 {
			if s[0] != '.' {
				return p, false
			}
			s = s[1:]
		}
		n, c, ok := dtoi(s)
		if !ok || n > 0xFF {
			return p, false
		}
		if c > 1 && s[0] == '0' {
			// Reject non-zero components with leading zeroes.
			return p, false
		}
		s = s[c:]
		p[i] = byte(n)
	}
	if len(s) != 0 {
		return p, false
	}
	return p, true
}

// Bigger than we need, not too big to worry about overflow
const big = 0xFFFFFF

// Decimal to integer.
// Returns number, characters consumed, success.
func dtoi(s string) (n int, i int, ok bool) {
	n = 0
	for i = 0; i < len(s) && '0' <= s[i] && s[i] <= '9'; i++ {
		n = n*10 + int(s[i]-'0')
		if n >= big {
			return big, i, false
		}
	}
	if i == 0 {
		return 0, 0, false
	}
	return n, i, true
}

// ParseIPV4 is shorthand for net.ParseIP(ip).To4()
func ParseIPV4(ip string) (IPV4, error) {
	p4, ok := parseIPV4(ip)
	if !ok {
		return p4, errors.Errorf("invalid IP address: %s", ip)
	}

	return p4, nil
}

func newIPV4(seg1, seg2, seg3, seg4 byte) IPV4 {
	//return IntToBytes(seg1, seg2, seg3, seg4)
	return IPV4{seg1, seg2, seg3, seg4}
}

const MaxIPv4StringLen = len("255.255.255.255")

func (p4 IPV4) String() string {
	b := make([]byte, MaxIPv4StringLen)

	n := ubtoa(b, 0, p4[0])
	b[n] = '.'
	n++

	n += ubtoa(b, n, p4[1])
	b[n] = '.'
	n++

	n += ubtoa(b, n, p4[2])
	b[n] = '.'
	n++

	n += ubtoa(b, n, p4[3])
	return string(b[:n])
}

// ubtoa encodes the string form of the integer v to dst[start:] and
// returns the number of bytes written to dst. The caller must ensure
// that dst has sufficient length.
func ubtoa(dst []byte, start int, v byte) int {
	if v < 10 {
		dst[start] = v + '0'
		return 1
	} else if v < 100 {
		dst[start+1] = v%10 + '0'
		dst[start] = v/10 + '0'
		return 2
	}

	dst[start+2] = v%10 + '0'
	dst[start+1] = (v/10)%10 + '0'
	dst[start] = v/100 + '0'
	return 3
}

const ByteLimit = byte(255)

// ip 自增 1
func (p4 IPV4) Add() IPV4 {
	for i := 3; i > 0; i-- {
		p4[i] += 1
		if p4[i] <= ByteLimit {
			break
		}
	}
	return p4
}

// 比较两个IPV4地址大小
// 结果：0 相等， -1 小于，1 大于
func (p4 IPV4) Compare(dest IPV4) int {
	var v int
	for i := range p4 {
		v = compareUint8(p4[i], dest[i])
		if v != 0 {
			return v
		}
	}
	return v
}

// 返回两个byte的比较结果
// 结果：0 相等， -1 小于，1 大于
func compareUint8(src, dest uint8) int {
	if src > dest {
		return 1
	}
	if src < dest {
		return -1
	}
	return 0
}
