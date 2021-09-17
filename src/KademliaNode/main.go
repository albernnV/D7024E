package main

import "github.com/albernnV/D7024E/pkg/kademlia"

func main() {
	net := &kademlia.Network{}
	//kademlia.Listen("", 8000)
	net.SendPingMessage()
}
