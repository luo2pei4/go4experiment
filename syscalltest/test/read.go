package test

import (
	"fmt"
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
