package main

import (
	"github.com/albernnV/D7024E/pkg/cli"
	"github.com/albernnV/D7024E/pkg/kademlia"
)

func main() {
	net := &kademlia.Network{}
	routingTable := &kademlia.RoutingTable{}
	
	go cli.Cli(net)
	go net.Listen()
}
