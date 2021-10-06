package kademlia

import (
	"bufio"
	"encoding/hex"
	"fmt"
	"net"
)

type Network struct {
	shortlistCh   chan []Contact //channel where shortlists from the goroutines will be written to
	inactiveNodes ContactCandidates
	routingTable  *RoutingTable
}

func (network *Network) Listen() {
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
		if err != nil {
			fmt.Println("Error reading from UDP stream")
		} else {
			incomingMessage := hex.EncodeToString(p)
			messageType := getTypeFromMessage(incomingMessage)
			switch messageType {
			case "FIND_NODE_RPC":
				targetContactAsString := incomingMessage[len(messageType)+1 : len(incomingMessage)-1]
				targetContact := StringToContact(targetContactAsString)
				closestContacts := network.routingTable.FindClosestContacts(targetContact.ID, bucketSize)
				senderID := NewKademliaID(incomingMessage[len(messageType)+1+len(targetContactAsString)+1:])
				sender := NewContact(senderID, remoteaddr.String())
				shortlistAsString := shortlistToString(&closestContacts)
				//Send shortlist to sender
				fmt.Fprintf(ser, shortlistAsString)
				network.routingTable.routingTableChan <- sender
			case "FIND_VALUE_RPC":
				//TODO: Lookup and return the value that's sought after
			case "STORE_VALUE_RPC":
				//TODO: Store value somewhere
			case "SHORTLIST":
				//TODO: Append shortlist to old shortlist
			case "PING":
				go sendPongResponse(ser, remoteaddr)
			case "PONG":
				//TODO: Update k-buckets
			}
		}
		/*fmt.Printf("Message receive from %v %s \n", remoteaddr, p)
		if err != nil {
			fmt.Printf("Some error  %v", err)
			continue
		}
		go sendResponse(ser, remoteaddr)*/
		//TODO: Check what the response is and send to a different channel depending on response
		//Things to check:
		//	- Shortlist
		//	- Value from a STORE_RPC
		//	- Ping message
		//	- FIND_NODE_RPC from another node
		//	- FIND_VALUE from another node
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

func sendPongResponse(conn *net.UDPConn, addr *net.UDPAddr) {
	_, err := conn.WriteToUDP([]byte("PONG "), addr)
	if err != nil {
		fmt.Printf("Couldn't send response %v", err)
	}
}

func (network *Network) SendFindContactMessage(contact *Contact, target *Contact) {
	//Establish connection
	conn, err := net.Dial("tcp", contact.Address)
	//Send contact to channel to mark as inactive
	if err != nil {
		network.shortlistCh <- []Contact{}
		network.inactiveNodes.Append([]Contact{*contact})
	}
	reader := bufio.NewReader(conn)
	targetAsString := target.String()
	//Send find node rpc together with the target contact
	fmt.Fprintf(conn, "FIND_NODE_RPC;"+targetAsString+";"+network.routingTable.me.ID.String()+"\n")
	//Wait for showrtlist as answer
	shortListString, err := reader.ReadString('\n')
	if err != nil {
		network.shortlistCh <- []Contact{}
		network.inactiveNodes.Append([]Contact{*contact})
	}
	shortList := preprocessShortlist(shortListString)
	network.shortlistCh <- shortList
}

//When receiving a shortlist it will be a string with the structure that looks like this:
//	"contact(ID, IP);contact(ID, IP)..."
//This function will convert this string into a list containing all the contacts
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

func shortlistToString(shortlist *[]Contact) string {
	var shortlistString string
	for _, contact := range *shortlist {
		contactString := contact.String()
		shortlistString = shortlistString + contactString + ";"
	}
	shortlistString = shortlistString[:len(shortlistString)-1] //remove last semicolon
	return shortlistString
}

//Takes as input a string structured as "contact(ID, IP) and converts it into a Contact"
func StringToContact(contactAsString string) Contact {
	var address string
	var id string
	contactRune := []rune(contactAsString)
	hasReadAddress := false
	//Skip 8 first letters since they always start with "contact("
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

func (network *Network) SendFindDataMessage(ID string, contact *Contact) {
	conn, err := net.Dial("tcp", contact.Address)
	if err != nil {
		fmt.Printf("Some error %v\n", err)
	}
	fmt.Fprintf(conn, "FIND_VALUE_RPC;"+ID+"\n")

	/** A function call to Listen() is needed here but Listen()
	needs to be redone bc that should be the only function that listens **/

}

func (network *Network) SendStoreMessage(data []byte, contact *Contact, target Contact) {
	conn, err := net.Dial("tcp", contact.Address)
	if err != nil {
		fmt.Printf("Some error %v\n", err)
	}

	dataToString := hex.EncodeToString(data)
	fmt.Fprintf(conn, "STORE_VALUE_RPC;"+dataToString+";"+target.ID.String())

	/** A function call to Listen() is needed here but Listen()
	needs to be redone bc that should be the only function that listens **/
}

func (network *Network) sendShortlist(recipient *Contact, shortlist *[]Contact) {
	shortlistAsString := shortlistToString(shortlist)

}

func getTypeFromMessage(message string) string {
	var messageType string
	for _, letter := range message {
		if string(letter) == ";" {
			break
		} else {
			messageType = messageType + string(letter)
		}
	}
	return messageType
}
