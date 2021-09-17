package kademlia

import (
	"bufio"
	"fmt"
	"net"
)

type Network struct {
	Name string
}

func Listen() {
	/*fmt.Println("Listening..........")
	address := ip + ":" + strconv.Itoa(port)
	ln, _ := net.Listen("tcp4", address)

	conn, _ := ln.Accept()
	message, _ := bufio.NewReader(conn).ReadString('\n')
	fmt.Println("Message recieved:" + string(message))
	if string(message) == "PING" {
		fmt.Fprintf(conn, "PONG")
	}*/
	fmt.Println("Listening.....")
	p := make([]byte, 2048)
	addr := net.UDPAddr{
		Port: 1234,
		IP:   net.ParseIP("0.0.0.0"),
	}
	ser, err := net.ListenUDP("udp", &addr)
	if err != nil {
		fmt.Printf("Some error %v\n", err)
		return
	}
	for {
		_, remoteaddr, err := ser.ReadFromUDP(p)
		fmt.Printf("Message receive from %v %s \n", remoteaddr, p)
		if err != nil {
			fmt.Printf("Some error  %v", err)
			continue
		}
		go sendResponse(ser, remoteaddr)
	}

}

func (network *Network) SendPingMessage() {
	/*conn, _ := net.Dial("tcp", address)
	fmt.Fprintf(conn, "PING"+"\n")
	message, _ := bufio.NewReader(conn).ReadString('\n')
	fmt.Println("Message recieved:" + string(message))
	if string(message) == "PONG" {
		fmt.Println("PONG recieved")
	}*/

	p := make([]byte, 2048)
	conn, err := net.Dial("udp", "172.18.0.3:1234")
	if err != nil {
		fmt.Printf("Some error %v", err)
		return
	}
	fmt.Fprintf(conn, "PIIING")
	_, err = bufio.NewReader(conn).Read(p)
	if err == nil {
		fmt.Printf("%s\n", p)
	} else {
		fmt.Printf("Some error %v\n", err)
	}
	conn.Close()

}

func sendResponse(conn *net.UDPConn, addr *net.UDPAddr) {
	_, err := conn.WriteToUDP([]byte("POONG "), addr)
	if err != nil {
		fmt.Printf("Couldn't send response %v", err)
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
