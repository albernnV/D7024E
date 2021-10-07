package main

import (
	"github.com/albernnV/D7024E/pkg/cli"
	"github.com/albernnV/D7024E/pkg/kademlia"
)

func main() {
	// id??? vart genereras och läggs in?
	alpha := 3
	contact := &kademlia.Contact{}
	node := kademlia.NewKademliaInstance(alpha, *contact)
	go cli.Cli(node)
	node.Start()
}
