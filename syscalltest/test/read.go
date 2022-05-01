package test

import (
	"fmt"
	"strings"
	"syscall"
)

type NetLinker struct {
	fd int
}

func NewNetLinker(groups, proto int) (*NetLinker, error) {

	s, err := syscall.Socket(syscall.AF_NETLINK, syscall.SOCK_RAW, proto)
	if err != nil {
		return nil, fmt.Errorf("socket: %s", err)
	}

	saddr := &syscall.SockaddrNetlink{
		Family: syscall.AF_NETLINK,
		Pid:    uint32(0),
		Groups: uint32(groups),
	}

	err = syscall.Bind(s, saddr)
	if err != nil {
		return nil, fmt.Errorf("bind: %s", err)
	}

	return &NetLinker{fd: s}, nil
}

func (l *NetLinker) RecvUdevEvent(knlMsgChan chan []byte) {
	for {
		pkt := make([]byte, 2048)
		_, err := syscall.Read(l.fd, pkt)
		if err != nil {
			fmt.Printf("syscall.Read error: %s", err)
		}
		knlMsgChan <- pkt
	}
}

type DevInfo struct {
	isUsbStorage  bool
	diskName      string
	partitionName string
}

type UsbMap map[string]*DevInfo

func (usbMap UsbMap) CheckUSBStorage(itemsMap map[string]string) {

	action := itemsMap["ACTION"]
	devName := itemsMap["DEVNAME"]
	devPath := itemsMap["DEVPATH"]
	devType := itemsMap["DEVTYPE"]
	subSystem := itemsMap["SUBSYSTEM"]

	containsDevPath := func(dp string) (string, bool) {
		if len(usbMap) == 0 {
			return dp, false
		}
		for key := range usbMap {
			if strings.Contains(dp, key) {
				return key, true
			}
		}
		return dp, false
	}

	if action == "add" {
		switch subSystem {
		case "usb":
			if dp, ok := containsDevPath(devPath); !ok {
				usbMap[dp] = &DevInfo{
					isUsbStorage: false,
				}
				fmt.Printf("Insert usb device, devPath: %s\n", dp)
			}
		case "block":
			if dp, ok := containsDevPath(devPath); ok {
				if devType == "disk" {
					usbMap[dp].diskName = devName
					fmt.Printf("usb device is a disk, diskName: %s\n", devName)
				}
				if devType == "partition" {
					usbMap[dp].isUsbStorage = true
					usbMap[dp].partitionName = devName
					fmt.Printf("usb device is a disk, partition: %s\n", devName)
				}
			}
		}
	} else if action == "remove" {
		if subSystem == "usb" {
			if dp, ok := containsDevPath(devPath); ok {
				delete(usbMap, dp)
				fmt.Printf("usb device %s pulled out.\n", dp)
			}
		}
	}
}
