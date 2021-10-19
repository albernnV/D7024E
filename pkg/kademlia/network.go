package kademlia

import (
	"fmt"
	"net"
	"strings"
)

type Network struct {
	shortlistCh      chan []Contact //channel where shortlists from the goroutines will be written to
	inactiveNodes    ContactCandidates
	routingTable     *RoutingTable
	storedValues     map[KademliaID]string
	closeNetwork     bool
	listenConnection *net.UDPConn
}

func NewNetwork(me Contact) *Network {
	shortlistCh := make(chan []Contact)
	inactiveNodes := ContactCandidates{[]Contact{}}
	routingTable := NewRoutingTable(me)
	storedValues := make(map[KademliaID]string)

	fmt.Println("Listening.....")
	addr := net.UDPAddr{
		Port: 8000,
		IP:   net.ParseIP(""),
	}
	listenConnection, err := net.ListenUDP("udp", &addr)
	if err != nil {
		fmt.Printf("Some error %v\n", err)
	}

	return &Network{shortlistCh, inactiveNodes, routingTable, storedValues, false, listenConnection}
}

func (network *Network) LoopListen() {
	for {
		network.Listen()
	}
}

// Listen acts as a global listener for incoming messages and processes each message depending on the type of the message
func (network *Network) Listen() {
	p := make([]byte, 2048)
	_, remoteaddr, err := network.listenConnection.ReadFromUDP(p)
	remoteaddr.Port = 8000
	if err != nil {
		fmt.Println(err)
	} else {
		incomingMessage := string(p)
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
			network.listenConnection.WriteToUDP([]byte("SHORTLIST;"+shortlistAsString+";"+network.routingTable.me.ID.String()+"\n"), remoteaddr)
			// Update buckets
			go network.routingTable.AddContact(sender)
		case "FIND_VALUE_RPC": //Lookup and return the value that's sought after
			IDAsString := data
			valueID := NewKademliaID(IDAsString)
			value := network.storedValues[*valueID]
			//Send value to sender
			network.listenConnection.WriteToUDP([]byte("VALUE;"+value+";"+network.routingTable.me.ID.String()+"\n"), remoteaddr)
			// Update buckets
			network.routingTable.AddContact(sender)
		case "STORE_VALUE_RPC":
			fmt.Println("store: " + data)
			valueID := HashingData([]byte(data))
			network.storedValues[*valueID] = data
			network.routingTable.AddContact(sender)
		case "SHORTLIST":
			shortlistAsString := data
			newShortlist := preprocessShortlist(shortlistAsString)
			go network.addToShortlist(newShortlist)
			network.routingTable.AddContact(sender)
		case "VALUE":
			value := data
			fmt.Println(value)
			go network.routingTable.AddContact(sender)
		case "PING":
			network.routingTable.AddContact(sender)
			network.listenConnection.WriteToUDP([]byte("PONG;0;"+network.routingTable.me.ID.String()+"\n"), remoteaddr)
		case "PONG":
			network.routingTable.AddContact(sender)
		}
	}
}

// Takes a message and returns message type, message data and sender ID
func preprocessIncomingMessage(message string) (string, string, string) {
	var messageType string
	var data string
	var senderID string
	numberOfEncounteredSemicolon := 0
	for _, letter := range message {
		if string(letter) == ";" {
			numberOfEncounteredSemicolon += 1
			continue
		} else if string(letter) == "\n" { //All messages end with a newline character
			break
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
	if shortlistString == "0" {
		return []Contact{}
	}
	var contactString string
	shortlist := make([]Contact, 0)
	for _, letter := range shortlistString {
		if string(letter) == ")" {
			newContact := StringToContact(contactString + string(letter))
			shortlist = append(shortlist, newContact)
			contactString = ""
		} else {
			contactString = contactString + string(letter)
		}
	}
	//Add last contact to shortlist
	//newContact := StringToContact(contactString)
	//shortlist = append(shortlist, newContact)
	return shortlist
}

//Turns a list of contacts into a string with the format
//	contact("ID", "IP")contact("ID", "IP")...
func shortlistToString(shortlist *[]Contact) string {
	if len(*shortlist) == 0 {
		return "0"
	}
	var shortlistString string
	for _, contact := range *shortlist {
		contactString := contact.String()
		shortlistString = shortlistString + contactString
	}
	return shortlistString
}

func (network *Network) addToShortlist(shortlist []Contact) {
	network.shortlistCh <- shortlist
}

//Takes as input a string structured as contact("ID", "IP", "Distance") and converts it into a Contact
// or contact("ID", "IP", "") if a distance hasn't been set
func StringToContact(contactAsString string) Contact {
	s := strings.Split(contactAsString, ", ") //Splits string into [contact("ID" "IP" "Distance")]
	ID := s[0]
	address := s[1]
	distance := s[2]

	ID = ID[9 : len(ID)-1]                //Remove contact(" and last "
	address = address[1 : len(address)-1] //Remove " from beginning and end of string
	newKademliaID := NewKademliaID(ID)
	newContact := NewContact(newKademliaID, address)

	if len(distance) > 3 { //A distance has been set for the contact
		distance = distance[1 : len(distance)-2] //Remove ") at the end of the string
		newContact.distance = NewKademliaID(distance)
	}
	return newContact
}

func (network *Network) SendPingMessage(contact *Contact) {
	conn, err := net.Dial("udp", contact.Address)
	if err != nil {
		fmt.Printf("Somee error %v", err)
		return
	}
	_, sendErr := fmt.Fprintf(conn, "PING;0;"+network.routingTable.me.ID.String()+"\n")
	network.Listen()
	if sendErr != nil {
		fmt.Printf("Some error %v\n", err)
	}
	conn.Close()
}

func (network *Network) sendPongResponse(conn *net.UDPConn, addr *net.UDPAddr) {
	_, err := conn.WriteToUDP([]byte("PONG;0;"+network.routingTable.me.ID.String()+"\n"), addr)
	if err != nil {
		fmt.Printf("Couldn't send response %v", err)
	}
}

func (network *Network) SendFindContactMessage(contact *Contact, target *Contact) {
	//Establish connection
	conn, err := net.Dial("udp", contact.Address)
	//Send contact to channel to mark as inactive
	if err != nil {
		network.shortlistCh <- []Contact{}
		network.inactiveNodes.Append([]Contact{*contact})
	}

	targetAsString := target.String()
	//Send find node rpc together with the target contact
	fmt.Fprintf(conn, "FIND_NODE_RPC;"+targetAsString+";"+network.routingTable.me.ID.String()+"\n")
	conn.Close()
}

func (network *Network) SendFindDataMessage(ID string, contact *Contact) {
	conn, err := net.Dial("udp", contact.Address)
	if err != nil {
		fmt.Printf("Some error %v\n", err)
	}
	fmt.Fprintf(conn, "FIND_VALUE_RPC;"+ID+";"+network.routingTable.me.ID.String()+"\n")
	conn.Close()
}

func (network *Network) SendStoreMessage(data []byte, contact *Contact) {
	conn, err := net.Dial("udp", contact.Address)
	if err != nil {
		fmt.Printf("Some error %v\n", err)
	}

	fmt.Fprintf(conn, "STORE_VALUE_RPC;"+string(data)+";"+network.routingTable.me.ID.String()+"\n")
	conn.Close()
}
