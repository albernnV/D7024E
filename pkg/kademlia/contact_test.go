package kademlia

import (
	"testing"
)

func TestNewContact(t *testing.T) {

	testID := NewKademliaID("0000000000000000000000000000000000000001")
	testAddr := "172.0.0.1:8000"

	newContact := NewContact(testID, testAddr)

	if newContact.ID != testID || newContact.Address != testAddr {
		t.Errorf("This does not return a new contact, got: %v and %v, want: %s and %s", newContact.ID, newContact.Address, testID, testAddr)
	}

}

func TestContactCalcDistance(t *testing.T) {

	kademliaID1 := NewKademliaID("0000000000000000000000000000000000000001")
	kademliaID2 := NewKademliaID("0000000000000000000000000000000000000002")
	testAddr := "172.0.0.1:8000"
	testDistance := NewKademliaID("0000000000000000000000000000000000000003")

	newContact := NewContact(kademliaID1, testAddr)

	newContact.CalcDistance(kademliaID2)

	if *newContact.distance != *testDistance {
		t.Errorf("This is not the same distance, got: %v, want: %v", newContact.distance, testDistance)
	}

}

func TestContactLess(t *testing.T) {

	kademliaID1 := NewKademliaID("0000000000000000000000000000000000000001")
	kademliaID2 := NewKademliaID("0000000000000000000000000000000000000002")

	contact1 := NewContact(kademliaID1, "")
	contact2 := NewContact(kademliaID2, "")

	contact1.distance = NewKademliaID("0000000000000000000000000000000000000003")
	contact2.distance = NewKademliaID("0000000000000000000000000000000000000004")

	newContact := contact1.distance.Less(contact2.distance)

	if newContact != true {
		t.Errorf("The distance is not smaller, got: %v, want: %v", newContact, contact1.distance)
	}
}

func TestContactString(t *testing.T) {

	kademliaID := NewKademliaID("0000000000000000000000000000000000000001")
	addr := "172.0.0.1:8000"

	newContact := NewContact(kademliaID, addr)
	contactToString := newContact.String()

	testString := `contact("0000000000000000000000000000000000000001", "172.0.0.1:8000")`

	if contactToString != testString {
		t.Errorf("This does not output the same string, got: %s, want: %s", contactToString, testString)
	}

}
