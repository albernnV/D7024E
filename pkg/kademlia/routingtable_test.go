package kademlia

import (
	"fmt"
	"testing"
)

func TestRoutingTable(t *testing.T) {
	rt := NewRoutingTable(NewContact(NewKademliaID("FFFFFFFF00000000000000000000000000000000"), "localhost:8000"))

	rt.AddContact(NewContact(NewKademliaID("FFFFFFFF00000000000000000000000000000000"), "localhost:8001"))
	rt.AddContact(NewContact(NewKademliaID("1111111100000000000000000000000000000000"), "localhost:8002"))
	rt.AddContact(NewContact(NewKademliaID("1111111200000000000000000000000000000000"), "localhost:8002"))
	rt.AddContact(NewContact(NewKademliaID("1111111300000000000000000000000000000000"), "localhost:8002"))
	rt.AddContact(NewContact(NewKademliaID("1111111400000000000000000000000000000000"), "localhost:8002"))
	rt.AddContact(NewContact(NewKademliaID("2111111400000000000000000000000000000000"), "localhost:8002"))

	contacts := rt.FindClosestContacts(NewKademliaID("2111111400000000000000000000000000000000"), 20)
	for i := range contacts {
		fmt.Println(contacts[i].String())
	}
}


func TestGetBucketIndex(t *testing.T) {
	me := NewContact(NewKademliaID("0000000000000000000000000000000000000000"), "")
	rt := NewRoutingTable(me)
	contactID := NewKademliaID("0000000000000000000000000000000000000000")
	bucketIndex := rt.getBucketIndex(contactID)
	// 1 is one byte so 0001
	correctIndex := 159
	if !(bucketIndex == correctIndex) {
		t.Errorf("Wrong bucket id, want: %v, got: %v", correctIndex, bucketIndex)
	}
}


func TestAddContactRT(t *testing.T) {
	me := NewContact(NewKademliaID("0000000000000000000000000000000000000000"), "")
	rt := NewRoutingTable(me)
	//should be in bucket with index 2: [2^2 - 2^3)
	contact := NewContact(NewKademliaID("1000000000000000000000000000000000000000"), "")
	rt.AddContact(contact)
	bucket := rt.buckets[3]
	bool := true
	for e := bucket.list.Front(); e != nil; e = e.Next() {
		nodeID := e.Value.(Contact).ID
		if (contact).ID.Equals(nodeID) {
			bool = false
		}
	}
	if bool {
		t.Errorf("The contact is not in the right bucket")
	}
}


func TestFindClosestContacts(t *testing.T) {
	me := NewContact(NewKademliaID("0000000000000000000000000000000000000000"), "")
	rt := NewRoutingTable(me)
	contact1 := NewContact(NewKademliaID("1000000000000000000000000000000000000000"), "")
	contact2 := NewContact(NewKademliaID("1100000000000000000000000000000000000000"), "")
	contact3 := NewContact(NewKademliaID("1110000000000000000000000000000000000000"), "")
	rt.AddContact(contact1)
	rt.AddContact(contact2)
	rt.AddContact(contact3)
	numContacts := 2

	closestLstTest := rt.FindClosestContacts(NewKademliaID("1110000000000000000000000000000000000000"), numContacts)
	var closestLstCorrect []Contact
	closestLstCorrect = append(closestLstCorrect, contact3)
	closestLstCorrect = append(closestLstCorrect, contact2)

	if len(closestLstTest) == numContacts {
		for i := 0; i < numContacts; i++ {
			if closestLstTest[i].ID != closestLstCorrect[i].ID {
				t.Errorf("Contacts is wrong, expected: %v, got: %v", closestLstCorrect[i].ID, closestLstTest[i].ID)
			}
		}
	} else {
		t.Errorf("Number of contacts is wrong, expected: %v, got: %v", numContacts, len(closestLstTest))
	}
}

