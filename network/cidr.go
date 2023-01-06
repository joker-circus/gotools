package network

import (
	"fmt"
	"strconv"
	"strings"
)

type IPV4Net struct {
	IP   IPV4
	Mask uint8
}

// return ipv4, masklen, error
func ParseIPV4Net(cidr string) (IPV4Net, error) {
	var info IPV4Net
	var err error
	message := fmt.Errorf("invalid IP address: %s", cidr)

	address := strings.Split(cidr, "/")
	if len(address) != 2 {
		return info, message
	}

	ip := address[0]
	info.IP, err = ParseIPV4(ip)
	if err != nil {
		return info, message
	}

	maskLen, err := strconv.Atoi(address[1])
	if err != nil {
		return info, message
	}
	info.Mask = uint8(maskLen)
	return info, nil
}

// 判断网段是否包含该IP
func (i IPV4Net) Contains(ip IPV4) bool {
	minIP, maxIP := i.getIpMaskRange()

	if ip.Compare(minIP) < 0 {
		return false
	}

	if ip.Compare(maxIP) > 0 {
		return false
	}
	return true
}

func (i IPV4Net) String() string {
	return fmt.Sprintf("%s/%d", i.IP.String(), i.Mask)
}

// 遍历主机地址，参数 f 用于接收各主机地址值，若返回 false 则结束
func (i IPV4Net) Range(f func(value IPV4) bool) {
	minIP, maxIP := i.getIpMaskRange()
	for minIP.Compare(maxIP) < 1 {
		if !f(minIP) {
			break
		}
		minIP = minIP.Add()
	}
}

// 最大主机数量
func (i IPV4Net) Len() uint {
	return GetCidrHostNum(i.Mask)
}

// 获取网段可分配地址范围，建议使用 Range 方法，防止溢出
func (i IPV4Net) GetNetworkIPS() []IPV4 {
	data := make([]IPV4, 0, i.Len())

	i.Range(func(value IPV4) bool {
		data = append(data, value)
		return true
	})
	return data
}

// 获取网段最大值和最小值
// return minIP, maxIP
func (i IPV4Net) getIpMaskRange() (IPV4, IPV4) {
	p4, maskLen := i.IP, i.Mask
	seg1MinIp, seg1MaxIp := getIpSeg1Range(p4[0], maskLen)
	seg2MinIp, seg2MaxIp := getIpSeg2Range(p4[1], maskLen)
	seg3MinIp, seg3MaxIp := getIpSeg3Range(p4[2], maskLen)
	seg4MinIp, seg4MaxIp := getIpSeg4Range(p4[3], maskLen)

	return newIPV4(seg1MinIp, seg2MinIp, seg3MinIp, seg4MinIp), newIPV4(seg1MaxIp, seg2MaxIp, seg3MaxIp, seg4MaxIp)
}

// 获取网段最大值和最小值
// 如：10.0.0.0/8，返回：10.0.0.1，10.255.255.254
// return minIP, maxIP
func GetCidrIpRange(cidr string) (IPV4, IPV4, error) {
	ipv4Net, err := ParseIPV4Net(cidr)
	if err != nil {
		return IPV4{}, IPV4{}, err
	}

	minIP, maxIP := ipv4Net.getIpMaskRange()
	return minIP, maxIP, nil
}

//得到第一段IP的区间（第一片段.第二片段.第三片段.第四片段）
func getIpSeg1Range(ipSeg, maskLen uint8) (uint8, uint8) {
	if maskLen > 8 {
		return ipSeg, ipSeg
	}
	return getIpSegRange(ipSeg, 8-maskLen)
}

//得到第二段IP的区间（第一片段.第二片段.第三片段.第四片段）
func getIpSeg2Range(ipSeg, maskLen uint8) (uint8, uint8) {
	if maskLen > 16 {
		return ipSeg, ipSeg
	}
	return getIpSegRange(ipSeg, 16-maskLen)
}

//得到第三段IP的区间（第一片段.第二片段.第三片段.第四片段）
func getIpSeg3Range(ipSeg, maskLen uint8) (uint8, uint8) {
	if maskLen > 24 {
		return ipSeg, ipSeg
	}
	return getIpSegRange(ipSeg, uint8(24-maskLen))
}

//得到第四段IP的区间（第一片段.第二片段.第三片段.第四片段）
func getIpSeg4Range(ipSeg, maskLen uint8) (uint8, uint8) {
	segMinIp, segMaxIp := getIpSegRange(ipSeg, uint8(32-maskLen))
	return segMinIp + 1, segMaxIp - 1
}

//根据用户输入的基础IP地址和CIDR掩码计算一个IP片段的区间
// 	192.168.1.0/6
// 	第一段参数 192 8-6    out  192,195
// 	第二段参数 168 16-6   out  0,255
// 	第一段参数 1 32-6     out  0,255
// 	第一段参数 0 32-6     out  0,255
func getIpSegRange(userSegIp, offset uint8) (uint8, uint8) {
	var ipSegMax uint8 = 255
	netSegIp := ipSegMax << offset
	segMinIp := netSegIp & userSegIp
	segMaxIp := userSegIp&(255<<offset) | ^(255 << offset)
	return segMinIp, segMaxIp
}

//计算得到CIDR地址范围内可拥有的主机数量
func GetCidrHostNum(maskLen uint8) uint {
	cidrIpNum := uint(0)
	var i uint = uint(32 - maskLen - 1)
	for ; i >= 1; i-- {
		cidrIpNum += 1 << i
	}
	return cidrIpNum
}
