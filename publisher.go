package main

import (
	//standard lib
	"encoding/json"
	"fmt"

	//external
	zmq "github.com/pebbe/zmq4"
)

func publisher(caps chan *DNSCapture, port string) {
	publisher, _ := zmq.NewSocket(zmq.PUB)
	connect_string := fmt.Sprintf("tcp://*:%s", port)
	publisher.Bind(connect_string)

	for cap := range caps {
		cap_json, err := json.Marshal(cap)
		if err != nil {
			panic(err)
		}
		fmt.Printf("%s\n", cap_json)
		msg := fmt.Sprintf("dns %s", cap_json)
		publisher.SendMessage(msg)
	}
}
