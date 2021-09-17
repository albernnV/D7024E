package kademlia

import (
	"bufio"
	"fmt"
	"net"
	"strconv"
)

type Network struct {
	Name string
}

func Listen(ip string, port int) {
	fmt.Println("Listening...")
	address := ip + ":" + strconv.Itoa(port)
	ln, _ := net.Listen("tcp", address)

	conn, _ := ln.Accept()
	message, _ := bufio.NewReader(conn).ReadString('\n')
	fmt.Println("Message recieved:" + string(message))
	if string(message) == "PING" {
		fmt.Fprintf(conn, "PONG")
	}
}

func (network *Network) SendPingMessage(address string) {
	conn, _ := net.Dial("tcp", address)
	fmt.Fprintf(conn, "PING"+"\n")
	message, err := bufio.NewReader(conn).ReadString('\n')
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Message recieved:" + string(message))
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
