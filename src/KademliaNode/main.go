package main

import (
	"github.com/albernnV/D7024E/pkg/cli"
	"github.com/albernnV/D7024E/pkg/kademlia"
)

func main() {
	go cli.Cli()
	kademliaID1 := kademlia.NewKademliaID("0000000000000000000000000000000000000001")
	me := kademlia.NewContact(kademliaID1, "")
	alpha := 3
	kademliaInstance := kademlia.NewKademliaInstance(alpha, me)
	kademliaInstance.Start()
}
