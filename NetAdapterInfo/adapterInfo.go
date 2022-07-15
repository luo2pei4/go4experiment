package main

import (
	"encoding/json"
	"fmt"
	"net"
	"strings"
)

type AdapterInfo struct {
	Index int      `json:"index"`
	Speed int      `json:"speed"`
	Name  string   `json:"name"`
	Ipv4  string   `json:"ipv4"`
	Ipv6  string   `json:"ipv6"`
	Mac   string   `json:"mac"`
	Flags []string `json:"flags"`
}

func GetAdaptersInfo() ([]AdapterInfo, error) {
	ifis, err := net.Interfaces()
	if err != nil {
		return nil, err
	}

	adapters := make([]AdapterInfo, 0, len(ifis))

	for _, ifi := range ifis {
		adapter := AdapterInfo{
			Index: ifi.Index,
			Speed: ifi.MTU,
			Name:  ifi.Name,
			Mac:   ifi.HardwareAddr.String(),
		}
		flags := strings.ToUpper(ifi.Flags.String())
		adapter.Flags = strings.Split(flags, "|")
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
		arr, err := json.Marshal(adapter)
		if err != nil {
			fmt.Println(err.Error())
		}
		fmt.Printf("%s\n", string(arr))
	}
}
