package main

import (
	"fmt"
	"go4experiment/syscalltest/test"
	"sync"
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

	var wg sync.WaitGroup
	wg.Add(1)
	go func(linker *test.NetLinker, wg *sync.WaitGroup) {
		for {
			// 此处调用RecvUdevEvent会产生阻塞
			msg, err := linker.RecvUdevEvent()
			if err != nil {
				fmt.Println(err.Error())
				wg.Done()
			} else {
				fmt.Println(msg)
			}
		}
	}(udev, &wg)
	wg.Wait()
}
