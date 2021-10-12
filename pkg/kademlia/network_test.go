package kademlia

import (
	"net"
	"testing"
)

var stringID string = "0000000000000000000000000000000000000001"
var kademliaID = NewKademliaID(stringID)
var me Contact = NewContact(kademliaID, "127.0.0.1:8000")

var network *Network = NewNetwork(me)

func TestStringToContact(t *testing.T) {
	contactID := "0000000000000000000000000000000000000002"
	contactAddress := "173.16.0.1"
	contactPlaceholder := `contact("0000000000000000000000000000000000000002", "173.16.0.1")`
	contact := StringToContact(contactPlaceholder)
	if contact.ID.String() != contactID || contact.Address != contactAddress {
		t.Errorf("Convertion to contact incorrect, got: (ID: %s, IP: %s) want: (ID: %s, IP: %s)", contact.ID.String(), contact.Address, contactID, contactAddress)
	}
}

func TestShortlistToString(t *testing.T) {
	contactID := NewKademliaID("0000000000000000000000000000000000000002")
	contactAddress := "173.16.0.1"
	node1 := NewContact(contactID, contactAddress)
	shortlist := []Contact{me, node1}
	shortlistString := shortlistToString(&shortlist)
	correctFormat := me.String() + ";" + node1.String()
	if shortlistString != correctFormat {
		t.Errorf("Convertion to string incorrect, got: %s want: %s", shortlistString, correctFormat)
	}

	emptyShortlist := []Contact{}
	emptyShortlistString := shortlistToString(&emptyShortlist)
	if emptyShortlistString != "0" {
		t.Errorf("Does not handle empty shortlist correctly got: %s want: %s", emptyShortlistString, "0")
	}
}

func TestPreprocessShortlist(t *testing.T) {
	contactID := NewKademliaID("0000000000000000000000000000000000000002")
	contactAddress := "173.16.0.1"
	node1 := NewContact(contactID, contactAddress)
	shortlist := []Contact{me, node1}
	shortlistString := shortlistToString(&shortlist)
	stringToShortlist := preprocessShortlist(shortlistString)
	if len(stringToShortlist) != len(shortlist) {
		t.Errorf("Contact list doesn't have the correct length, got: %d want: %d", len(stringToShortlist), len(shortlist))
	}
	for i := 0; i < len(shortlist); i++ {
		if stringToShortlist[i].ID.String() != shortlist[i].ID.String() || stringToShortlist[i].Address != shortlist[i].Address {
			t.Errorf("Incorrect conversion from string to shortlist, got: %v want: %v", stringToShortlist, shortlist)
		}
	}

	//Test empty case
	emptyShortlist := preprocessShortlist("0")
	if len(emptyShortlist) != 0 {
		t.Errorf("Empty shortlist string not handled correctly, got: %v want: %v", emptyShortlist, []Contact{})
	}
}

func TestPreprocessIncomingMessage(t *testing.T) {
	message := "FIND_VALUE_RPC;00001;00002"
	messageType, data, senderIDString := preprocessIncomingMessage(message)
	if messageType != "FIND_VALUE_RPC" || data != "00001" || senderIDString != "00002" {
		t.Errorf(
			"Preprocessing of message is incorrect, got: (message type: %s, data: %s, sender: %s) want: (message type: %s, data: %s, sender: %s)",
			messageType,
			data,
			senderIDString,
			"FIND_VALUE_RPC",
			"00001",
			"00002",
		)
	}
}

func TestSendPingMessage(t *testing.T) {
	p := make([]byte, 2048)
	addr := net.UDPAddr{
		Port: 8000,
		IP:   net.ParseIP(""),
	}
	go network.SendPingMessage(&me)
	conn, _ := net.ListenUDP("udp", &addr)
	conn.ReadFromUDP(p)
	incomingMessage := string(p)
	messageType, _, _ := preprocessIncomingMessage(incomingMessage)
	if messageType != "PING" {
		t.Errorf("PING message was not sent")
	}
	conn.Close()
}

func TestSendPongResponse(t *testing.T) {
	p := make([]byte, 2048)
	addr := net.UDPAddr{
		Port: 8000,
		IP:   net.ParseIP(""),
	}
	homeAddr := net.UDPAddr{
		Port: 8000,
		IP:   net.ParseIP("127.0.0.1"),
	}
	conn, _ := net.ListenUDP("udp", &addr)
	go network.sendPongResponse(conn, &homeAddr)
	conn.ReadFromUDP(p)
	messageType, _, _ := preprocessIncomingMessage(string(p))
	if messageType != "PONG" {
		t.Errorf("PONG message was not sent")
	}
	conn.Close()
}

func TestSendFindContactMessage(t *testing.T) {
	contactID := NewKademliaID("0000000000000000000000000000000000000002")
	contactAddress := "127.0.0.1:8000"
	node1 := NewContact(contactID, contactAddress)

	p := make([]byte, 2048)
	addr := net.UDPAddr{
		Port: 8000,
		IP:   net.ParseIP(""),
	}
	go network.SendFindContactMessage(&node1, &me)
	conn, _ := net.ListenUDP("udp", &addr)
	conn.ReadFromUDP(p)
	messageType, data, senderID := preprocessIncomingMessage(string(p))
	if messageType != "FIND_NODE_RPC" {
		t.Errorf("FIND_NODE_RPC message was not sent")
	}
	if data != me.String() || senderID != me.ID.String() {
		t.Errorf("FIND_NODE_RPC message did not contain correct information, got: %s want: %s", string(p), "FIND_NODE_RPC;"+me.String()+";"+me.ID.String())
	}
	conn.Close()
}

func TestSendFindDataMessage(t *testing.T) {
	contactID := NewKademliaID("0000000000000000000000000000000000000002")
	contactAddress := "127.0.0.1:8000"
	node1 := NewContact(contactID, contactAddress)

	p := make([]byte, 2048)
	addr := net.UDPAddr{
		Port: 8000,
		IP:   net.ParseIP(""),
	}
	dataID := "0000010520300000050000000067800000000002"
	go network.SendFindDataMessage(dataID, &node1)
	conn, _ := net.ListenUDP("udp", &addr)
	conn.ReadFromUDP(p)
	messageType, data, senderID := preprocessIncomingMessage(string(p))
	if messageType != "FIND_VALUE_RPC" {
		t.Errorf("FIND_VALUE_RPC message was not sent")
	}
	if data != dataID || senderID != me.ID.String() {
		t.Errorf("FIND_VALUE_RPC message did not contain correct information, got: %s want: %s", string(p), "FIND_VALUE_RPC;"+dataID+";"+me.ID.String())
	}
	conn.Close()
}

func TestSendStoreMessage(t *testing.T) {
	dataToStore := "Hej jag heter Albernn"

	contactID := NewKademliaID("0000000000000000000000000000000000000002")
	contactAddress := "127.0.0.1:8000"
	node1 := NewContact(contactID, contactAddress)

	p := make([]byte, 2048)
	addr := net.UDPAddr{
		Port: 8000,
		IP:   net.ParseIP(""),
	}
	go network.SendStoreMessage([]byte(dataToStore), &node1)
	conn, _ := net.ListenUDP("udp", &addr)
	conn.ReadFromUDP(p)
	messageType, data, senderID := preprocessIncomingMessage(string(p))
	if messageType != "STORE_VALUE_RPC" {
		t.Errorf("STORE_VALUE_RPC message was not sent")
	}
	if data != dataToStore || senderID != me.ID.String() {
		t.Errorf("STORE_VALUE_RPC message did not contain correct information, got: %s want: %s", string(p), "STORE_VALUE_RPC;"+dataToStore+";"+me.ID.String())
	}
	conn.Close()
}

/*func TestListen(t *testing.T) {
	p := make([]byte, 2048)
	contactID := NewKademliaID("0000000000000000000000000000000000000002")
	contactAddress := "127.0.0.1:8000"
	node1 := NewContact(contactID, contactAddress)
	go network.Listen()
	conn, _ := net.Dial("udp", node1.Address)
	fmt.Fprintf(conn, "FIND_NODE_RPC;"+me.String()+";"+me.ID.String()+"\n")
	conn.Read(p)
	messageType, _, _ := preprocessIncomingMessage(string(p))
	if messageType != "SHORTLIST" {
		t.Errorf("SHORTLIST message not sent")
	}
}*/
