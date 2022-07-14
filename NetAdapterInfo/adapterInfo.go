package main

import (
	"fmt"
	"net"
	"strings"
)

type AdapterInfo struct {
	Index      int
	MTU        int
	Name       string
	Ipv4       string
	Ipv6       string
	MacAddress string
	Flags      []string
}

func GetAdaptersInfo() ([]AdapterInfo, error) {
	ifis, err := net.Interfaces()
	if err != nil {
		return nil, err
	}

	adapters := make([]AdapterInfo, 0, len(ifis))

	for _, ifi := range ifis {
		adapter := AdapterInfo{
			Index:      ifi.Index,
			MTU:        ifi.MTU,
			Name:       ifi.Name,
			MacAddress: ifi.HardwareAddr.String(),
			Flags:      strings.Split(ifi.Flags.String(), "|"),
		}
		addrs, err := ifi.Addrs()
		if err != nil {
			return nil, err
		}
		for idx, addr := range addrs {
			ip, _, err := net.ParseCIDR(addr.String())
			if err != nil {
				return nil, err
			}
			if idx%2 == 0 {
				adapter.Ipv4 = ip.String()
			} else {
				adapter.Ipv6 = ip.String()
			}
		}
		adapters = append(adapters, adapter)
	}

	return adapters, nil
}

func main() {
	adapters, err := GetAdaptersInfo()
	if err != nil {
		fmt.Println(err.Error())
	}

	for _, adapter := range adapters {
		fmt.Printf("[%d](%d)%s\t%v\t%s\n", adapter.Index, adapter.MTU, adapter.Name, adapter.Flags, adapter.MacAddress)
		fmt.Printf("[IPv4/IPv6]%s/%s\n", adapter.Ipv4, adapter.Ipv6)
	}
}
