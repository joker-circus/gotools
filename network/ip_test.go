package network

import (
	"fmt"
	"testing"
)

func TestGetNetworkIPS(t *testing.T) {
	ipv4Net, err := ParseIPV4Net("192.168.0.0/8")
	if err != nil {
		panic(err)
	}
	data := ipv4Net.GetNetworkIPS()
	fmt.Println(len(data))

	if len(data) == 0 {
		return
	}
	fmt.Println(data[0], data[len(data) - 1])
}
