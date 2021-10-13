package main

import (
	"github.com/albernnV/D7024E/pkg/cli"
	"github.com/albernnV/D7024E/pkg/kademlia"
)

func main() {
	// id??? vart genereras och läggs in?
	alpha := 3
	me := &kademlia.Contact{}
	kademliaInstance := kademlia.NewKademliaInstance(alpha, *me)
	go cli.Cli(kademliaInstance)
	kademliaInstance.Start()
}
