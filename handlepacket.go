package main

import (
	//standard lib
	"time"

	//external
	"code.google.com/p/gopacket"
	"code.google.com/p/gopacket/layers"
)

func handlePacket(packet gopacket.Packet) (*DNSCapture, error) {
	dcap := new(DNSCapture)
	var err error

	if dnslayer := packet.Layer(layers.LayerTypeDNS); dnslayer != nil {

		iplayer := packet.Layer(layers.LayerTypeIPv4)
		if iplayer == nil {
			iplayer = packet.Layer(layers.LayerTypeIPv6)
			if iplayer == nil {
				panic("Could not find IPv4 or IPv6")
			}
			ip, _ := iplayer.(*layers.IPv6)
			dcap.SrcIP = ip.SrcIP
			dcap.DstIP = ip.DstIP

		} else {
			ip, _ := iplayer.(*layers.IPv4)
			dcap.SrcIP = ip.SrcIP
			dcap.DstIP = ip.DstIP
		}

		dns, _ := dnslayer.(*layers.DNS)

		dcap.Query = string(dns.Questions[0].Name)
		dcap.Timestamp = time.Now().Unix()

	}
	return dcap, err

}
