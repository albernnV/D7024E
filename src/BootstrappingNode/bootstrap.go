package main

import (
	"github.com/albernnV/D7024E/pkg/cli"
	"github.com/albernnV/D7024E/pkg/kademlia"
)

func main() {
	go cli.Cli()
	net := &kademlia.Network{}
	for {
		kademlia.Listen()
	}
}
