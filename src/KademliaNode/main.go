package main

import (
	"github.com/albernnV/D7024E/pkg/cli"
	"github.com/albernnV/D7024E/pkg/kademlia"
)

func main() {
	net := &kademlia.Network{}
	go cli.Cli(net)
	//kademlia.Listen("", 8000)
	go net.Listen()
}

