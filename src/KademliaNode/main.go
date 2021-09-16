package main

import "github.com/albernnV/D7024E/src/pkg/kademlia"

func main() {
	net := &kademlia.Network{}
	//kademlia.Listen("", 8000)
	net.SendPingMessage("172.18.0.3:8000")
}
