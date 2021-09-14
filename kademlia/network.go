package kademlia

import (
	"fmt"
)

type Network struct {
}

func Listen(ip string, port int) {
	fmt.Println("Listening to " + ip)
}

func (network *Network) SendPingMessage() {
	fmt.Println("Sending message...")
}

func (network *Network) SendFindContactMessage(contact *Contact) {
	// TODO
}

func (network *Network) SendFindDataMessage(hash string) {
	// TODO
}

func (network *Network) SendStoreMessage(data []byte) {
	// TODO
}
