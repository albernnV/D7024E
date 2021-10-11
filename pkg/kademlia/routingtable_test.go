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

func TestgetBucketIndex(t *testing.T) {
	me := NewContact(NewKademliaID("0000000000000000000000000000000000000000"), "")
	rt := NewRoutingTable(me)
	contactID := NewKademliaID("0000000000000000000000000000000000000001")
	bucketIndex := rt.getBucketIndex(contactID)
	correctIndex := 99999999
	fmt.Println("index: "+ string(bucketIndex))
	if !(bucketIndex == correctIndex) {
		t.Errorf("Wrong bucket id, want: %v, got: %v", correctIndex, bucketIndex)
	}

}

/*
func TestAddContactRT(t *testing.T) {
	me := NewContact(NewKademliaID("0000000000000000000000000000000000000000"), "")
	rt := NewRoutingTable(me)
	//should be in bucket with index 2: [2^2 - 2^3)
	contact := NewContact(NewKademliaID("0000000000000000000000000000000000000001"), "")
	rt.AddContact(contact)
	bucket := rt.buckets[0]
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
*/

/*
func TestFindClosestContacts(t *testing.T) {

}

func TestUpdateRoutingTable(t *testing.T) {

}
*/