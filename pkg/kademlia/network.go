package kademlia

import (
	"fmt"
	"net"
	"strings"
)

// The network file is used to send various messages to different nodes in the network.
// There are seven different message types that can be sent which are:
// - FIND_NODE_RPC: Used to find a target node and nodes closest to the target in the network
// - FIND_VALUE_RPC: Used to retreive data from a node
// - STORE_VALUE_RPC: Used to store data at a node
// - SHORTLIST: Used to send a list of contacts to a node, often as a response to a FIND_NODE_RPC
// - VALUE: Used to send data to a node
// - PING: Used to ping a node
// - PONG: Used as a response when a PING is received
//
// Messages are structured like MESSAGE_TYPE;DATA;SENDER_ID where
// - MESSAGE_TYPE is equal to one of the seven message types described above
// - DATA is equal to the data that is sent with the message
// - SENDER_ID is the ID of the sender node
//
// DATA that should be in a message depending on MESSAGE_TYPE
// - FIND_NODE_RPC: DATA = contact as a string
// - FIND_VALUE_RPC: DATA = hashed value of the string data that needs to be retreived
// - STORE_VALUE_RPC: DATA = string to stored
// - SHORTLIST: DATA = shortlist as a string
// - VALUE: DATA = string of data that was requested
// - PING: DATA = 0
// - PONG: DATA = 0

type Network struct {
	shortlistCh       chan []Contact //channel where shortlists from the goroutines will be written to
	inactiveNodes     ContactCandidates
	routingTable      *RoutingTable
	storedValues      map[KademliaID]string // Used to store data that are received in a STORE_VALUE_RPC
	receivedValue     ValueAndSender        // Used to temporary store data that has been requested by a FIND_VALUE_RPC
	receviedValueChan chan ValueAndSender   // Temporarily store data that has been received from a VALUE message
	listenConnection  *net.UDPConn
}

type ValueAndSender struct {
	value  string
	sender Contact
}

// Creates a network to able to use it as a reference
func NewNetwork(me Contact) *Network {
	shortlistCh := make(chan []Contact)
	receivedValueChan := make(chan ValueAndSender)
	inactiveNodes := ContactCandidates{[]Contact{}}
	routingTable := NewRoutingTable(me)
	storedValues := make(map[KademliaID]string)

	addr := net.UDPAddr{
		Port: 8000,
		IP:   net.ParseIP(""),
	}
	listenConnection, err := net.ListenUDP("udp", &addr)
	if err != nil {
		fmt.Printf("Some error %v\n", err)
	}

	return &Network{shortlistCh, inactiveNodes, routingTable, storedValues, ValueAndSender{}, receivedValueChan, listenConnection}
}

func (network *Network) LoopListen() {
	for {
		network.Listen()
	}
}

// Listen acts as a global listener for incoming messages and processes each message depending on the type of the message
func (network *Network) Listen() {
	p := make([]byte, 3048)
	_, remoteaddr, err := network.listenConnection.ReadFromUDP(p)
	remoteaddr.Port = 8000
	if err != nil {
		fmt.Println(err)
	} else {
		incomingMessage := string(p)
		network.handleMessage(incomingMessage, remoteaddr)
	}
}

// Handles a message and preforms various functions depending in the message type
// Every time a message is received the routing table is updated with the sender contact
func (network *Network) handleMessage(message string, remoteaddr *net.UDPAddr) {
	messageType, data, senderIDAsString := preprocessIncomingMessage(message)
	senderID := NewKademliaID(senderIDAsString)
	sender := NewContact(senderID, remoteaddr.String())
	switch messageType {
	case "FIND_NODE_RPC": // Send sortlist of nodes closest to target back to sender
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
	case "STORE_VALUE_RPC": // Store the data that was accopanied with the message
		valueID := HashingData([]byte(data))
		network.storedValues[*valueID] = data
		network.routingTable.AddContact(sender)
	case "SHORTLIST": // Append the received shortlist to the current shortlist
		shortlistAsString := data
		newShortlist := preprocessShortlist(shortlistAsString)
		go network.addToShortlist(newShortlist)
		network.routingTable.AddContact(sender)
	case "VALUE": // Send value to the channel to be processed by another goroutine
		value := data
		network.receivedValue.value = value
		network.receivedValue.sender = sender
		network.receviedValueChan <- network.receivedValue
		go network.routingTable.AddContact(sender)
	case "PING": //Send a PONG response message
		network.routingTable.AddContact(sender)
		network.listenConnection.WriteToUDP([]byte("PONG;0;"+network.routingTable.me.ID.String()+"\n"), remoteaddr)
	case "PONG":
		network.routingTable.AddContact(sender)
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
//	"contact(ID, IP, distance)contact(ID, IP, distance)..."
//This function will convert this string into a list containing all the contacts in the string
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
	return shortlist
}

//Turns a list of contacts into a string with the format
//	contact("ID", "IP", "distance")contact("ID", "IP", "distance")...
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

// Sends a PING message to see if a contact is online and waits for a PONG response
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

// Sends a FIND_NODE_RPC to a contact
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

// Sends a FIND_VALUE_RPC to a contact
func (network *Network) SendFindDataMessage(ID string, contact Contact) {
	conn, err := net.Dial("udp", contact.Address)
	if err != nil {
		fmt.Printf("Some error %v\n", err)
	}
	fmt.Fprintf(conn, "FIND_VALUE_RPC;"+ID+";"+network.routingTable.me.ID.String()+"\n")
	conn.Close()
}

// Sends a STORE_VALUE_RPC to a contact
func (network *Network) SendStoreMessage(data []byte, contact *Contact) {
	conn, err := net.Dial("udp", contact.Address)
	if err != nil {
		fmt.Printf("Some error %v\n", err)
	}
	fmt.Fprintf(conn, "STORE_VALUE_RPC;"+string(data)+";"+network.routingTable.me.ID.String()+"\n")
	conn.Close()
}
