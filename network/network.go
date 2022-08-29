package network

import (
	"net"
)

const (
	UNKNOWN_IP_ADDR = "-"
)

var localIP string

func init() {
	localIP = getLocalIp()
}

// LocalIP returns host's ip
func LocalIP() string {
	return localIP
}

// getLocalIp enumerates local net interfaces to find local ip, it should only be called in init phase
func getLocalIp() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return UNKNOWN_IP_ADDR
	}

	for _, addr := range addrs {
		if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			//ipv4 := ipnet.IP.To4()
			//if ipv4 == nil {
			//	continue
			//}
			return ipnet.IP.String()
		}
	}
	return UNKNOWN_IP_ADDR
}

// 返回域名对应的 Host，和 Hosts 方法类似，不过仅返回一个值。
func Host(domain string) (string, error) {
	ip, err := net.ResolveIPAddr("ip", domain)
	if err != nil {
		return "", err
	}

	return ip.String(), nil
}

// 返回域名对应的 Hosts 列表
func Hosts(domain string) (addrs []string, err error)  {
	return net.LookupHost(domain)
}
