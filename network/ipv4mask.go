package network

import (
	"fmt"
	"net"
	"strconv"
	"strings"
)

// 解析 cidr 如：192.168.255.45/25
// 返回 192.168.255.45，25
func ParseCIDR(cidr string) (string, int, error) {
	ip, ipNet, err := net.ParseCIDR(cidr)
	if err != nil {
		return "", 0, err
	}
	ones, _ := ipNet.Mask.Size()
	return ip.String(), ones, nil
}

// ContainsCIDR 子网a 是否包含 子网b
// b 是 a 的子集
// return true - b是a的子网; false b 不是 a 的子网
func ContainsCIDR(a, b *net.IPNet) bool {
	ones1, _ := a.Mask.Size()
	ones2, _ := b.Mask.Size()
	return ones1 <= ones2 && a.Contains(b.IP)
}

// 将整型数字掩码转换为 ip 格式
// 如 24 对应为 255.255.255.0
func MaskToString(ones int) (string, error) {
	m := net.CIDRMask(ones, 8 * net.IPv4len)
	if len(m) != 4 {
		return "", fmt.Errorf("ipv4Mask: len must be 4 bytes")
	}

	return fmt.Sprintf("%d.%d.%d.%d", m[0], m[1], m[2], m[3]), nil
}

// 将ip格式的掩码转换为整型数字
// 如 255.255.255.0 对应的整型数字为 24
func MaskToInt(netmask string) (int, error) {
	ipSplitArr := strings.Split(netmask, ".")
	if len(ipSplitArr) != 4 {
		return 0, fmt.Errorf("netmask:%v is not valid, pattern should like: 255.255.255.0", netmask)
	}
	ipv4MaskArr := make([]byte, 4)
	for i, value := range ipSplitArr {
		intValue, err := strconv.Atoi(value)
		if err != nil {
			return 0, fmt.Errorf("ipMaskToInt call strconv.Atoi error:[%v] string value is: [%s]", err, value)
		}
		if intValue > 255 {
			return 0, fmt.Errorf("netmask cannot greater than 255, current value is: [%s]", value)
		}
		ipv4MaskArr[i] = byte(intValue)
	}

	ones, _ := net.IPv4Mask(ipv4MaskArr[0], ipv4MaskArr[1], ipv4MaskArr[2], ipv4MaskArr[3]).Size()
	return ones, nil
}
