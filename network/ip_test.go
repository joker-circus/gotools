package network

import (
	"fmt"
	"testing"
)

func TestGetNetworkIPS(t *testing.T) {
	ipv4Net, err := ParseIPV4Net("192.168.0.0/24")
	if err != nil {
		panic(err)
	}
	data := ipv4Net.GetNetworkIPS()
	fmt.Println(len(data))

	if len(data) == 0 {
		return
	}
	fmt.Println(data[0], data[len(data)-1])
}

func BenchmarkParseIPV4(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_, _ = ParseIPV4("192.168.255.254")
	}
}

func BenchmarkIPV4_String(b *testing.B) {
	ip, _ := ParseIPV4("192.168.255.254")
	b.StartTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = ip.String()
	}
}

func BenchmarkIPV4Net_GetNetworkIPS(b *testing.B) {
	ipv4Net, err := ParseIPV4Net("192.168.0.0/24")
	if err != nil {
		panic(err)
	}
	b.StartTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = ipv4Net.GetNetworkIPS()
	}
}
