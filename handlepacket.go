package main

import (
	//standard lib
	"net"
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

func getIPAddress(packet gopacket.Packet) (net.IP, net.IP) {
	var srcip, dstip net.IP

	iplayer := packet.Layer(layers.LayerTypeIPv4)
	if iplayer == nil {
		iplayer = packet.Layer(layers.LayerTypeIPv6)
		if iplayer == nil {
			panic("Could not find IPv4 or IPv6")
		}
		ip, _ := iplayer.(*layers.IPv6)
		srcip = ip.SrcIP
		dstip = ip.DstIP

	} else {
		ip, _ := iplayer.(*layers.IPv4)
		srcip = ip.SrcIP
		dstip = ip.DstIP
	}

	return srcip, dstip

}

func getType(current layers.DNSType) string {
	var retval string
	if current == layers.DNSTypeA {
		retval = "A"
	} else if current == layers.DNSTypeAAAA {
		retval = "AAAA"
	} else if current == layers.DNSTypeNS {
		retval = "NS"
	} else if current == layers.DNSTypeMD {
		retval = "MD"
	} else if current == layers.DNSTypeMF {
		retval = "MF"
	} else if current == layers.DNSTypeCNAME {
		retval = "CNAME"
	} else if current == layers.DNSTypeSOA {
		retval = "SOA"
	} else if current == layers.DNSTypeMB {
		retval = "MB"
	} else if current == layers.DNSTypeMG {
		retval = "MG"
	} else if current == layers.DNSTypeMR {
		retval = "MR"
	} else if current == layers.DNSTypeNULL {
		retval = "NullRR"
	} else if current == layers.DNSTypeWKS {
		retval = "WKS"
	} else if current == layers.DNSTypePTR {
		retval = "PTR"
	} else if current == layers.DNSTypeHINFO {
		retval = "HINFO"
	} else if current == layers.DNSTypeMINFO {
		retval = "MINFO"
	} else if current == layers.DNSTypeMX {
		retval = "MX"
	} else if current == layers.DNSTypeTXT {
		retval = "TXT"
	} else if current == layers.DNSTypeSRV {
		retval = "SRV"
	}

	return retval
}

func processer(packets chan gopacket.Packet, captures chan interface{}) {
	for packet := range packets {
		if dnslayer := packet.Layer(layers.LayerTypeDNS); dnslayer != nil {
			dns, _ := dnslayer.(*layers.DNS)

			srcip, dstip := getIPAddress(packet)
			timestamp := time.Now().Unix()
			request := dns.QR

			if len(dns.Questions) > 0 {
				for _, query := range dns.Questions {
					dcap := new(DNSQuery)
					dcap.SrcIP = srcip
					dcap.DstIP = dstip
					dcap.Timestamp = timestamp
					dcap.Request = request
					dcap.Query = string(query.Name)
					dcap.Type = getType(query.Type)

					// write packet to wire
					captures <- dcap
				}
			}

			// add answers to packet
			if dns.QR {
				for _, answer := range dns.Answers {
					dcap := new(DNSAnswer)
					dcap.SrcIP = srcip
					dcap.DstIP = dstip
					dcap.Timestamp = timestamp
					dcap.Request = request
					dcap.Query = string(answer.Name)
					dcap.Type = getType(answer.Type)
					//if answer.Type == layers.DNSTypeA || answer.Type == layers.DNSTypeAAAA {
					dcap.Answer = answer.IP.String()
					dcap.TTL = answer.TTL
					//}

					// write packet to wire
					captures <- dcap
				}
			}

		}
	}
}
