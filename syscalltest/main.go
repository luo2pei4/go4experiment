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

	for msg := range knlMsgChan {
		str := strings.ReplaceAll(string(msg), "\x00", " ")
		fmt.Println(strings.TrimSpace(str))
	}
}
