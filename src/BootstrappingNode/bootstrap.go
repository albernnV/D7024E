package main

import (
	"github.com/albernnV/D7024E/pkg/cli"
	"github.com/albernnV/D7024E/pkg/kademlia"
)

func main() {
	alpha := 3
	contact := &kademlia.Contact{}
	node := kademlia.NewKademliaInstance(alpha, *contact)
	go cli.Cli(node)
	node.Start()
}
