package main

import (
	//standard lib
	"net"
)

type DNSCapture struct {
	SrcIP, DstIP net.IP
	Query        string
	Timestamp    int64
}
