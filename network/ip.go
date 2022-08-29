package network

import (
	"fmt"
	"strconv"
	"strings"
)

type IPV4 struct {
	Seg1 int
	Seg2 int
	Seg3 int
	Seg4 int
}

const IPLimit = 255

func ParseIPV4(ip string) (IPV4, error) {
	var info IPV4
	var err error
	message := fmt.Errorf("invalid IP address: %s", ip)

	ipSegs := strings.Split(ip, ".")
	if len(ipSegs) != 4 {
		return info, message
	}

	info.Seg1, err = strconv.Atoi(ipSegs[0])
	if err != nil {
		return info, message
	}
	info.Seg2, err = strconv.Atoi(ipSegs[1])
	if err != nil {
		return info, message
	}
	info.Seg3, err = strconv.Atoi(ipSegs[2])
	if err != nil {
		return info, message
	}
	info.Seg4, err = strconv.Atoi(ipSegs[3])
	if err != nil {
		return info, message
	}

	return info, nil
}

func newIPV4(seg1, seg2, seg3, seg4 int) IPV4 {
	return IPV4{seg1, seg2, seg3, seg4}
}

func (i IPV4) String() string {
	return fmt.Sprintf("%d.%d.%d.%d", i.Seg1, i.Seg2, i.Seg3, i.Seg4)
}

// ip 自增 1
func (i IPV4) Add() IPV4 {
	i.Seg4 += 1
	if i.Seg4 > IPLimit {
		i.Seg3 += 1
		i.Seg4 = i.Seg4 - IPLimit - 1
	}
	if i.Seg3 > IPLimit {
		i.Seg2 += 1
		i.Seg3 = i.Seg3 - IPLimit - 1
	}
	if i.Seg2 > IPLimit {
		i.Seg1 += 1
		i.Seg2 = i.Seg2 - IPLimit - 1
	}
	if i.Seg1 > IPLimit {
		i.Seg1 = i.Seg1 - IPLimit - 1
	}
	return i
}

// 比较两个IPV4地址大小
// 结果：0 相等， -1 小于，1 大于
func (i IPV4) Compare(dest IPV4) int {
	v := compareInt(i.Seg1, dest.Seg1)
	if v != 0 {
		return v
	}
	v = compareInt(i.Seg2, dest.Seg2)
	if v != 0 {
		return v
	}
	v = compareInt(i.Seg3, dest.Seg3)
	if v != 0 {
		return v
	}
	return compareInt(i.Seg4, dest.Seg4)
}

// 返回两个int的比较结果
// 结果：0 相等， -1 小于，1 大于
func compareInt(src, dest int) int {
	if src > dest {
		return 1
	}
	if src < dest {
		return -1
	}
	return 0
}
