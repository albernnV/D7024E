package main

import "github.com/albernnV/D7024E/pkg/kademlia"

func main() {
	//net := &kademlia.Network{}
	for {
		kademlia.Listen()
	}
}
