package main

import (
	//standard lib
	"time"

	//external
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
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
		dcap.Request = dns.QR

		// add answers to packet
		if dns.QR {
			for _, answer := range dns.Answers {
				if answer.Type == layers.DNSTypeA || answer.Type == layers.DNSTypeAAAA {
					dcap.Answers = append(dcap.Answers, answer.IP)
				}
			}
		}

	}
	return dcap, err

}

func processer(packets chan gopacket.Packet, captures chan *DNSCapture) {
	for packet := range packets {
		dcap, err := handlePacket(packet)
		if err != nil {
			panic(err)
		}
		captures <- dcap
	}
}
