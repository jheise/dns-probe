package main

import (

	//external
	"code.google.com/p/gopacket"
	"code.google.com/p/gopacket/pcap"
)

func main() {
	caps := make(chan *DNSCapture)

	// start capture consumer
	go publisher(caps)

	handle, err := pcap.OpenLive("wlan0", 1600, true, 0)
	if err != nil {
		panic(err)
	}
	err = handle.SetBPFFilter("udp and port 53")
	if err != nil {
		panic(err)
	}
	packetSource := gopacket.NewPacketSource(handle, handle.LinkType())
	for packet := range packetSource.Packets() {
		dcap, err := handlePacket(packet)
		if err != nil {
			panic(err)
		}
		caps <- dcap
	}
}
