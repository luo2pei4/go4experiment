package main

import (
	"encoding/json"
	"fmt"
	"strings"
	"syscall"
)

const (
	NUL = "\x00"
)

const (
	MsgTagAction    = "ACTION"
	MsgTagDevPath   = "DEVPATH"
	MsgTagDevTpye   = "DEVTYPE"
	MsgTagDevName   = "DEVNAME"
	MsgTagSubSystem = "SUBSYSTEM"
	MsgTagPartNum   = "PARTN"
	MsgTagPartName  = "PARTNAME"
	MsgTagSeqNum    = "SEQNUM"
)

type KObjUeventMsg struct {
	Action    string
	DevPath   string
	DevTpye   string
	DevName   string
	SubSystem string
	PartNum   string
	PartName  string
	SeqNum    string
}

func main() {

	fd, err := syscall.Socket(syscall.AF_NETLINK, syscall.SOCK_RAW, syscall.NETLINK_KOBJECT_UEVENT)
	if err != nil {
		fmt.Println("syscall.Socket", err.Error())
		return
	}

	socketAddrNetlink := &syscall.SockaddrNetlink{
		Family: syscall.AF_NETLINK,
		Groups: uint32(syscall.RTNLGRP_LINK),
	}

	if err = syscall.Bind(fd, socketAddrNetlink); err != nil {
		fmt.Println("syscall.Bind", err.Error())
		return
	}

	msg := make([]byte, 2048)
	for {
		n, err := syscall.Read(fd, msg)
		if err != nil {
			fmt.Println("syscall.Read", err.Error())
			continue
		}
		koum := parseKObjUeventMsg(msg[:n])
		if arr, err := json.Marshal(koum); err == nil {
			fmt.Println(string(arr))
		}
	}
}

func parseKObjUeventMsg(recMsg []byte) *KObjUeventMsg {
	items := strings.Split(string(recMsg), NUL)
	if len(items) > 1 && items[1][:6] != MsgTagAction {
		return nil
	}
	koum := &KObjUeventMsg{}
	for _, item := range items[1:] {
		if !strings.Contains(item, "=") {
			continue
		}
		kv := strings.Split(item, "=")
		switch kv[0] {
		case MsgTagAction:
			koum.Action = kv[1]
		case MsgTagDevPath:
			koum.DevPath = kv[1]
		case MsgTagDevTpye:
			koum.DevTpye = kv[1]
		case MsgTagDevName:
			koum.DevName = kv[1]
		case MsgTagSubSystem:
			koum.SubSystem = kv[1]
		case MsgTagPartNum:
			koum.PartNum = kv[1]
		case MsgTagPartName:
			koum.PartName = kv[1]
		case MsgTagSeqNum:
			koum.SeqNum = kv[1]
		}
	}
	return koum
}
