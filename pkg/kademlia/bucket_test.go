package kademlia

import (
	"testing"
)

func TestLen(t *testing.T) {
	bucket := newBucket()
	contact1 := NewContact(NewKademliaID("0000000000000000000000000000000000000001"), "")
	contact2 := NewContact(NewKademliaID("0000000000000000000000000000000000000002"), "")
	bucket.AddContact(contact1)
	bucket.AddContact(contact2)
	length := bucket.list.Len()
	if length != 2 {
		t.Errorf("This did not return correct length: got: %v, want: %v", length, 2)
	}
}

func TestAddContact(t *testing.T) {
	bucket := newBucket()
	contact1 := NewContact(NewKademliaID("0000000000000000000000000000000000000001"), "")
	contact2 := NewContact(NewKademliaID("0000000000000000000000000000000000000002"), "")

	// test adds to bucket
	bucket.AddContact(contact1)
	e := bucket.list.Front()
	nodeID := e.Value.(Contact).ID
	if bucket.list.Len() != 1 || nodeID != contact1.ID {
		t.Errorf("This did not add the new contact to the bucket")
	}

	//test adds to front for new node
	bucket.AddContact(contact2)
	e2 := bucket.list.Front()
	nodeID2 := e2.Value.(Contact).ID
	if bucket.list.Len() != 2 || nodeID2 != contact2.ID {
		t.Errorf("This did not add the new contact to the front of the bucket")
	}

	//test adds to front if already exist
	bucket.AddContact(contact1)
	e3 := bucket.list.Front()
	nodeID3 := e3.Value.(Contact).ID
	if bucket.list.Len() != 2 || nodeID3 != contact1.ID {
		t.Errorf("This did not move existing contact to the front")
	}
}

func TestGetContactAndCalcDistance(t *testing.T) {
	bucket := newBucket()
	contact1 := NewContact(NewKademliaID("0000000000000000000000000000000000000001"), "")
	contact2 := NewContact(NewKademliaID("0000000000000000000000000000000000000002"), "")
	targetContactID := NewKademliaID("0000000000000000000000000000000000000003")
	bucket.AddContact(contact1)
	bucket.AddContact(contact2)
	contactArray := bucket.GetContactAndCalcDistance(targetContactID)

	var contactArrayWanted []Contact
	contactArray1 := Contact{NewKademliaID("0000000000000000000000000000000000000001"), "", NewKademliaID("0000000000000000000000000000000000000002")}
	contactArray2 := Contact{NewKademliaID("0000000000000000000000000000000000000002"), "", NewKademliaID("0000000000000000000000000000000000000001")}
	contactArrayWanted = append(contactArrayWanted, contactArray1)
	contactArrayWanted = append(contactArrayWanted, contactArray2)

	bool := true
	count := len(contactArray)
	for i := 0; i < count; i++ {
		if contactArray[i] != contactArrayWanted[i] {
			bool = false
		}
	}

	if (len(contactArray) == len(contactArrayWanted)) && bool {
		t.Errorf("This did not give back the same contacts with the correct distances")
	}

}
