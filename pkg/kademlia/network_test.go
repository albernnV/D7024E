package kademlia

import (
	"testing"
)

var stringID string = "0000000000000000000000000000000000000001"
var kademliaID = NewKademliaID(stringID)
var me Contact = NewContact(kademliaID, "127.0.0.1:8000")

//var network *Network = NewNetwork(me)

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
}
