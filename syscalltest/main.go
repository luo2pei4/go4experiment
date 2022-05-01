package main

import (
	"fmt"
	"go4experiment/syscalltest/test"
	"strings"
	"syscall"
)

func main() {

	udev, err := test.NewNetLinker(
		syscall.RTNLGRP_LINK|syscall.RTNLGRP_IPV4_IFADDR|syscall.RTNLGRP_IPV6_IFADDR,
		syscall.NETLINK_KOBJECT_UEVENT,
	)
	if err != nil {
		fmt.Println(err.Error())
	}

	knlMsgChan := make(chan []byte)
	go udev.RecvUdevEvent(knlMsgChan)

	usbMap := make(test.UsbMap)

	for msg := range knlMsgChan {

		str := strings.ReplaceAll(string(msg), "\x00", " ")
		str = strings.TrimSpace(str)
		lines := strings.Split(str, "\n")
		for _, line := range lines {
			items := strings.Split(line, " ")
			itemsMap := make(map[string]string, len(items)-1)
			for _, item := range items {
				kv := strings.Split(item, "=")
				if len(kv) == 2 {
					itemsMap[kv[0]] = kv[1]
				}
			}
			usbMap.CheckUSBStorage(itemsMap)
		}

	}
}
