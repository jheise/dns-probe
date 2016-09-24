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

type DNSQuery struct {
	SrcIP, DstIP net.IP
	Request      bool
	Query        string
	Timestamp    int64
	Type         uint16
}

type DNSAnswer struct {
	SrcIP, DstIP, Answer net.IP
	Request              bool
	Timestamp            int64
	Type                 uint16
}
