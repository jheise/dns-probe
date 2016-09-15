package main

import (
	//standard lib
	"net"
)

type DNSCapture struct {
	SrcIP, DstIP net.IP
	Request      bool
	Query        string
	Answers      []net.IP
	Timestamp    int64
}
