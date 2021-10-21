package kademlia

import (
	"net"
	"testing"
)

var stringID string = "0000000000000000000000000000000000000001"
var dist string = "0000001000020000000000030000000000000001"
var kademliaID = NewKademliaID(stringID)
var me Contact = NewContact(kademliaID, "127.0.0.1:8000")

func TestStringToContact(t *testing.T) {
	contactID := "0000000000000000000000000000000000000002"
	contactDistance := "0000000000000000000000000000000000000005"
	contactAddress := "173.16.0.1"
	contactPlaceholder := `contact("0000000000000000000000000000000000000002", "173.16.0.1", "0000000000000000000000000000000000000005")`
	contact := StringToContact(contactPlaceholder)
	if contact.ID.String() != contactID || contact.Address != contactAddress || contact.distance.String() != contactDistance {
		t.Errorf("Convertion to contact incorrect, got: (ID: %s, IP: %s, dist: %s) want: (ID: %s, IP: %s, dist: %s)", contact.ID.String(), contact.Address, contact.distance.String(), contactID, contactAddress, contactDistance)
	}

	contactPlaceholder = `contact("0000000000000000000000000000000000000002", "173.16.0.1", "")`
	contact = StringToContact(contactPlaceholder)
	if contact.ID.String() != contactID || contact.Address != contactAddress || contact.distance != nil {
		t.Errorf("Convertion to contact incorrect, got: (ID: %s, IP: %s, dist: %s) want: (ID: %s, IP: %s, dist: %s)", contact.ID.String(), contact.Address, contact.distance.String(), contactID, contactAddress, contactDistance)
	}
}

func TestShortlistToString(t *testing.T) {
	me.distance = NewKademliaID(dist)

	contactID := NewKademliaID("0000000000000000000000000000000000000002")
	contactAddress := "173.16.0.1"
	contactDistance := "0000000000000000000000000000000000000005"
	node1 := NewContact(contactID, contactAddress)
	node1.distance = NewKademliaID(contactDistance)
	shortlist := []Contact{me, node1}
	shortlistString := shortlistToString(&shortlist)
	correctFormat := me.String() + node1.String()
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
	me.distance = NewKademliaID(dist)

	contactID := NewKademliaID("0000000000000000000000000000000000000002")
	contactAddress := "173.16.0.1"
	node1 := NewContact(contactID, contactAddress)
	node1.distance = NewKademliaID("0000000000000000000000000000000000000005")
	shortlist := []Contact{me, node1}
	shortlistString := shortlistToString(&shortlist)
	stringToShortlist := preprocessShortlist(shortlistString)
	if len(stringToShortlist) != len(shortlist) {
		t.Errorf("Contact list doesn't have the correct length, got: %d want: %d", len(stringToShortlist), len(shortlist))
	}
	for i := 0; i < len(shortlist); i++ {
		if stringToShortlist[i].ID.String() != shortlist[i].ID.String() || stringToShortlist[i].Address != shortlist[i].Address || stringToShortlist[i].distance.String() != shortlist[i].distance.String() {
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
	kademliaInstance.network.SendPingMessage(&me)
	kademliaInstance.network.listenConnection.ReadFromUDP(p)
	incomingMessage := string(p)
	messageType, _, _ := preprocessIncomingMessage(incomingMessage)
	if messageType != "PONG" {
		t.Errorf("PING message was not sent")
	}
}

func TestSendFindContactMessage(t *testing.T) {
	contactID := NewKademliaID("0000000000000000000000000000000000000002")
	contactAddress := "127.0.0.1:8000"
	node1 := NewContact(contactID, contactAddress)

	p := make([]byte, 2048)
	go kademliaInstance.network.SendFindContactMessage(&node1, &me)
	kademliaInstance.network.listenConnection.ReadFromUDP(p)
	messageType, data, senderID := preprocessIncomingMessage(string(p))
	if messageType != "FIND_NODE_RPC" {
		t.Errorf("FIND_NODE_RPC message was not sent")
	}
	if data != me.String() || senderID != me.ID.String() {
		t.Errorf("FIND_NODE_RPC message did not contain correct information, got: %s want: %s", string(p), "FIND_NODE_RPC;"+me.String()+";"+me.ID.String())
	}
}

func TestSendFindDataMessage(t *testing.T) {
	contactID := NewKademliaID("0000000000000000000000000000000000000002")
	contactAddress := "127.0.0.1:8000"
	node1 := NewContact(contactID, contactAddress)

	p := make([]byte, 2048)
	dataID := "0000010520300000050000000067800000000002"
	go kademliaInstance.network.SendFindDataMessage(dataID, node1)
	kademliaInstance.network.listenConnection.ReadFromUDP(p)
	messageType, data, senderID := preprocessIncomingMessage(string(p))
	if messageType != "FIND_VALUE_RPC" {
		t.Errorf("FIND_VALUE_RPC message was not sent")
	}
	if data != dataID || senderID != me.ID.String() {
		t.Errorf("FIND_VALUE_RPC message did not contain correct information, got: %s want: %s", string(p), "FIND_VALUE_RPC;"+dataID+";"+me.ID.String())
	}
}

func TestSendStoreMessage(t *testing.T) {
	dataToStore := "Hej jag heter Albernn"

	contactID := NewKademliaID("0000000000000000000000000000000000000002")
	contactAddress := "127.0.0.1:8000"
	node1 := NewContact(contactID, contactAddress)

	p := make([]byte, 2048)
	go kademliaInstance.network.SendStoreMessage([]byte(dataToStore), &node1)
	kademliaInstance.network.listenConnection.ReadFromUDP(p)
	messageType, data, senderID := preprocessIncomingMessage(string(p))
	if messageType != "STORE_VALUE_RPC" {
		t.Errorf("STORE_VALUE_RPC message was not sent")
	}
	if data != dataToStore || senderID != me.ID.String() {
		t.Errorf("STORE_VALUE_RPC message did not contain correct information, got: %s want: %s", string(p), "STORE_VALUE_RPC;"+dataToStore+";"+me.ID.String())
	}
}

func TestHandleMessageFindNode(t *testing.T) {
	remoteaddr := net.UDPAddr{
		Port: 8000,
		IP:   net.ParseIP("127.0.0.1"),
	}
	p := make([]byte, 2048)
	incomingMessage := "FIND_NODE_RPC;" + me.String() + ";0000000000000000000000000000000000000002"
	go kademliaInstance.network.handleMessage(incomingMessage, &remoteaddr)
	kademliaInstance.network.listenConnection.ReadFromUDP(p)
	messageType, _, _ := preprocessIncomingMessage(string(p))
	if messageType != "SHORTLIST" {
		t.Errorf("SHORTLIST message was not sent")
	}
}

func TestHandleMessageFindValue(t *testing.T) {
	hashedData := HashingData([]byte("hej jag heter Albernn"))
	remoteaddr := net.UDPAddr{
		Port: 8000,
		IP:   net.ParseIP("127.0.0.1"),
	}
	p := make([]byte, 2048)
	incomingMessage := "FIND_VALUE_RPC;" + hashedData.String() + ";0000000000000000000000000000000000000002"
	go kademliaInstance.network.handleMessage(incomingMessage, &remoteaddr)
	kademliaInstance.network.listenConnection.ReadFromUDP(p)
	messageType, _, _ := preprocessIncomingMessage(string(p))
	if messageType != "VALUE" {
		t.Errorf("SHORTLIST message was not sent")
	}
}

func TestHandleMessageStoreValue(t *testing.T) {
	dataToStore := "hej jag heter Albernn"
	dataHash := HashingData([]byte(dataToStore))
	remoteaddr := net.UDPAddr{
		Port: 8000,
		IP:   net.ParseIP("127.0.0.1"),
	}
	incomingMessage := "STORE_VALUE_RPC;" + dataToStore + ";0000000000000000000000000000000000000002"
	kademliaInstance.network.handleMessage(incomingMessage, &remoteaddr)
	if kademliaInstance.network.storedValues[*dataHash] != dataToStore {
		t.Errorf("Data was not stored")
	}
}

func TestHandleMessageShortlist(t *testing.T) {
	contactID := NewKademliaID("0000000000000000000000000000000000000002")
	contactAddress := "127.0.0.1:8000"
	node1 := NewContact(contactID, contactAddress)
	shortlist := []Contact{me, node1}
	shortlistString := shortlistToString(&shortlist)
	remoteaddr := net.UDPAddr{
		Port: 8000,
		IP:   net.ParseIP("127.0.0.1"),
	}
	incomingMessage := "SHORTLIST;" + shortlistString + ";0000000000000000000000000000000000000002"
	go kademliaInstance.network.handleMessage(incomingMessage, &remoteaddr)
	receivedShortlist := <-kademliaInstance.network.shortlistCh
	if receivedShortlist[0].ID.String() != shortlist[0].ID.String() || receivedShortlist[1].ID.String() != shortlist[1].ID.String() {
		t.Errorf("Shortlist was not received got: %v want: %v", receivedShortlist, shortlist)
	}
}

func TestHandleMessageValue(t *testing.T) {
	data := "hej jag heter Albernn"
	remoteaddr := net.UDPAddr{
		Port: 8000,
		IP:   net.ParseIP("127.0.0.1"),
	}
	incomingMessage := "VALUE;" + data + ";0000000000000000000000000000000000000002"
	go kademliaInstance.network.handleMessage(incomingMessage, &remoteaddr)
	receivedValue := <-kademliaInstance.network.receviedValueChan
	if receivedValue.value != data || receivedValue.sender.ID.String() != "0000000000000000000000000000000000000002" {
		t.Errorf("Correct value not received got: %s want: %s", receivedValue.value, data)
	}
}

func TestHandleMessagePing(t *testing.T) {
	remoteaddr := net.UDPAddr{
		Port: 8000,
		IP:   net.ParseIP("127.0.0.1"),
	}
	p := make([]byte, 2048)
	incomingMessage := "PING;0;0000000000000000000000000000000000000002"
	go kademliaInstance.network.handleMessage(incomingMessage, &remoteaddr)
	kademliaInstance.network.listenConnection.ReadFromUDP(p)
	messageType, _, _ := preprocessIncomingMessage(string(p))
	if messageType != "PONG" {
		t.Errorf("PONG message was not sent")
	}
}
