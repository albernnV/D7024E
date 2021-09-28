package kademlia

import (
	"bufio"
	"fmt"
	"net"
)

type Network struct {
}

func Listen() {
	fmt.Println("Listening.....")
	p := make([]byte, 2048)
	addr := net.UDPAddr{
		Port: 8000,
		IP:   net.ParseIP(""),
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
	p := make([]byte, 2048)
	conn, err := net.Dial("udp", "172.18.0.3:8000")
	if err != nil {
		fmt.Printf("Somee error %v", err)
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

func (network *Network) SendFindContactMessage(contact *Contact, target *Contact, hasNotAnsweredChannel chan Contact) []Contact {
	conn, err := net.Dial("tcp", contact.Address)
	if err != nil {
		hasNotAnsweredChannel <- *contact
		return []Contact{}
	}
	reader := bufio.NewReader(conn)
	targetAsString := target.String()
	fmt.Fprintf(conn, "FIND_NODE_RPC;"+targetAsString+"\n")
	shortListString, err := reader.ReadString('\n')
	shortList := preprocessShortlist(shortListString)
	if err != nil {
		hasNotAnsweredChannel <- *contact
		return []Contact{}
	}
	return shortList
}

func preprocessShortlist(shortlistString string) []Contact {
	var contactString string
	shortlist := make([]Contact, 0)
	for _, letter := range shortlistString {
		if string(letter) == ";" {
			newContact := StringToContact(contactString)
			shortlist = append(shortlist, newContact)
			contactString = ""
		} else {
			contactString = contactString + string(letter)
		}
	}
	return shortlist
}

func StringToContact(contactAsString string) Contact {
	var address string
	var id string
	contactRune := []rune(contactAsString)
	hasReadAddress := false
	for i := 8; i < len(contactRune); i++ {
		if string(contactRune[i]) == ")" {
			break
		}
		if string(contactRune[i]) == "," {
			hasReadAddress = true
			i = i + 2
		}
		if hasReadAddress {
			id = id + string(contactRune[i])
		} else {
			address = address + string(contactRune[i])
		}
	}
	newKademliaID := NewKademliaID(id)
	newContact := NewContact(newKademliaID, address)
	return newContact

}

func (network *Network) SendFindDataMessage(hash string) {
	// TODO
}

func (network *Network) SendStoreMessage(data []byte) {
	// TODO
}
