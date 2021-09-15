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
	address := ip + ":" + strconv.Itoa(port)
	ln, _ := net.Listen("tcp", address)

	conn, _ := ln.Accept()
	message, _ := bufio.NewReader(conn).ReadString('\n')
	if message == "PING" {
		fmt.Fprint(conn, "PONG")
	}
}

func (network *Network) SendPingMessage(address string) {
	conn, _ := net.Dial("tcp", address)
	fmt.Fprintf(conn, "PING"+"\n")
	message, _ := bufio.NewReader(conn).ReadString('\n')
	if message == "PONG" {
		fmt.Println("PONG recieved")
	}
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
