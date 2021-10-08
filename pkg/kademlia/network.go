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
	storedValues  map[KademliaID]string
}

func NewNetwork(me Contact) *Network {
	shortlistCh := make(chan []Contact)
	inactiveNodes := ContactCandidates{[]Contact{}}
	routingTable := NewRoutingTable(me)
	storedValues := make(map[KademliaID]string)
	return &Network{shortlistCh, inactiveNodes, routingTable, storedValues}
}

// Listen acts as a global listener for incoming messages and processes each message depending on the type of the message
func (network *Network) Listen() {
	fmt.Println("Listening.....")
	p := make([]byte, 2048)
	addr := net.UDPAddr{
		Port: 8000,
		IP:   net.ParseIP(""),
	}
	conn, err := net.ListenUDP("udp", &addr)
	if err != nil {
		fmt.Printf("Some error %v\n", err)
		return
	}
	for {
		_, remoteaddr, err := conn.ReadFromUDP(p)
		if err != nil {
			fmt.Println("Error reading from UDP stream")
		} else {
			incomingMessage := hex.EncodeToString(p)
			messageType, data, senderIDAsString := preprocessIncomingMessage(incomingMessage)
			senderID := NewKademliaID(senderIDAsString)
			sender := NewContact(senderID, remoteaddr.String())
			switch messageType {
			case "FIND_NODE_RPC":
				targetContactAsString := data
				targetContact := StringToContact(targetContactAsString)
				closestContacts := network.routingTable.FindClosestContacts(targetContact.ID, bucketSize)
				shortlistAsString := shortlistToString(&closestContacts)
				//Send shortlist to sender
				conn.WriteToUDP([]byte("SHORTLIST;"+shortlistAsString+";"+network.routingTable.me.ID.String()+"\n"), remoteaddr)
				// Update buckets
				go network.routingTable.AddContact(sender)
			case "FIND_VALUE_RPC": //Lookup and return the value that's sought after
				IDAsString := data
				valueID := NewKademliaID(IDAsString)
				value := network.storedValues[*valueID]
				//Send value to sender
				conn.WriteToUDP([]byte("VALUE;"+value+";"+network.routingTable.me.ID.String()+"\n"), remoteaddr)
				// Update buckets
				go network.routingTable.AddContact(sender)
			case "STORE_VALUE_RPC":
				valueID := HashingData([]byte(data))
				network.storedValues[*valueID] = data
				go network.routingTable.AddContact(sender)
			case "SHORTLIST":
				shortlistAsString := data
				newShortlist := preprocessShortlist(shortlistAsString)
				go network.addToShortlist(newShortlist)
				go network.routingTable.AddContact(sender)
			case "VALUE":
				value := data
				fmt.Println(value)
				go network.routingTable.AddContact(sender)
			case "PING":
				go network.routingTable.AddContact(sender)
				sendPongResponse(conn, remoteaddr)
			case "PONG":
				go network.routingTable.AddContact(sender)
			}
		}
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

// Takes a message and returns message type, message data and sender ID
func preprocessIncomingMessage(message string) (string, string, string) {
	var messageType string
	var data string
	var senderID string
	m := message[:len(message)-1] //Remove \n character at end
	numberOfEncounteredSemicolon := 0
	for _, letter := range m {
		if string(letter) == ";" {
			numberOfEncounteredSemicolon += 1
			continue
		}

		if numberOfEncounteredSemicolon == 0 {
			messageType = messageType + string(letter)
		} else if numberOfEncounteredSemicolon == 1 {
			data = data + string(letter)
		} else {
			senderID = senderID + string(letter)
		}
	}
	return messageType, data, senderID
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

//Turns a list of contacts into a string with the format
//	contact("ID", "IP");contact("ID", "IP")...
func shortlistToString(shortlist *[]Contact) string {
	var shortlistString string
	for _, contact := range *shortlist {
		contactString := contact.String()
		shortlistString = shortlistString + contactString + ";"
	}
	shortlistString = shortlistString[:len(shortlistString)-1] //remove last semicolon
	return shortlistString
}

func (network *Network) addToShortlist(shortlist []Contact) {
	network.shortlistCh <- shortlist
}

//Takes as input a string structured as contact("ID", "IP") and converts it into a Contact
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
	fmt.Fprintf(conn, "FIND_VALUE_RPC;"+ID+";"+network.routingTable.me.ID.String()+"\n")

	/** A function call to Listen() is needed here but Listen()
	needs to be redone bc that should be the only function that listens **/

}

func (network *Network) SendStoreMessage(data []byte, contact *Contact, target Contact) {
	conn, err := net.Dial("tcp", contact.Address)
	if err != nil {
		fmt.Printf("Some error %v\n", err)
	}

	dataToString := hex.EncodeToString(data)
	fmt.Fprintf(conn, "STORE_VALUE_RPC;"+dataToString+";"+network.routingTable.me.ID.String()) //TODO: Change so that it's sender id instead of target id since the target id can be hashed from the data

	/** A function call to Listen() is needed here but Listen()
	needs to be redone bc that should be the only function that listens **/
}
