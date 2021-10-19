package main

import (
	"log"
	"net"

	"github.com/albernnV/D7024E/pkg/cli"
	"github.com/albernnV/D7024E/pkg/kademlia"
)

func main() {
	alpha := 3
	ID := kademlia.NewRandomKademliaID()
	IPAddress := GetOutboundIP().String() + ":8000"
	me := kademlia.NewContact(ID, IPAddress)
	kademliaInstance := kademlia.NewKademliaInstance(alpha, me)
	kademliaInstance.Start()
	cli.Cli(kademliaInstance)
}

func GetOutboundIP() net.IP {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)

	return localAddr.IP
}
