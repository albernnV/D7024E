package main

import (
	"project/kademlia"
)

func main() {
	kademlia.Listen("127.0.0.1", 8000)
}
