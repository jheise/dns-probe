package main

import (

	//internal
	"fmt"

	//external
	"github.com/google/gopacket"
	"github.com/google/gopacket/pcap"
	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	iface = kingpin.Arg("interface", "interface to sniff traffic on").Default("wlan0").String()
	port  = kingpin.Arg("port", "port to expose for zmq").Default("7777").String()
)

func main() {
	kingpin.Parse()
	caps := make(chan *DNSCapture)
	fmt.Printf("Listening on %s\nPublishing on %s\n", *iface, *port)

	// start capture consumer
	go publisher(caps, *port)

	handle, err := pcap.OpenLive(*iface, 1600, true, 0)
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
