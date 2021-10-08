package kademlia

import "testing"

/*var stringID string = "0000000000000000000000000000000000000001"
var kademliaID = NewKademliaID(stringID)
var me Contact = NewContact(kademliaID, "127.0.0.1:8000")
var network *Network = NewNetwork(me)*/

func TestStringToContact(t *testing.T) {
	contactID := "0000000000000000000000000000000000000002"
	contactAddress := "173.16.0.1"
	contactPlaceholder := `contact("0000000000000000000000000000000000000002", "173.16.0.1")`
	contact := StringToContact(contactPlaceholder)
	if contact.ID.String() != contactID || contact.Address != contactAddress {
		t.Errorf("Convertion to contact incorrect, got: (ID: %s, IP: %s) want: (ID: %s, IP: %s)", contact.ID.String(), contact.Address, contactID, contactAddress)
	}
}
