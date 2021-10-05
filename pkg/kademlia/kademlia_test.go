package kademlia

import (
	"testing"
)

func TestFindClosestContact(t *testing.T) {
	//0000000000000000000000000000000000000002
	//0000000000000000000000000000000000000003 closest
	//0000000000000000000000000000000000000004

	//0000000000000000000000000000000000000001 TARGET
	kademliaID1 := NewKademliaID("0000000000000000000000000000000000000001")
	kademliaID2 := NewKademliaID("0000000000000000000000000000000000000002")
	kademliaID3 := NewKademliaID("0000000000000000000000000000000000000003")
	kademliaID4 := NewKademliaID("0000000000000000000000000000000000000004")
	node1 := NewContact(kademliaID1, "172.16.0.2:8000")
	node2 := NewContact(kademliaID2, "172.16.0.3:8000")
	node3 := NewContact(kademliaID3, "172.16.0.4:8000")
	node4 := NewContact(kademliaID4, "172.16.0.5:8000")
	contacts := []Contact{node2, node3, node4}

	closestContact := findClosestContact(contacts, &node1)
	if closestContact.ID != node3.ID {
		t.Errorf("Contact ID was incorrect, got %s, want: %s", closestContact.ID.String(), node3.ID.String())
	}
}

func TestSendFindNodeRPC(t *testing.T) {
	kademliaID1 := NewKademliaID("0000000000000000000000000000000000000001")
	kademliaID2 := NewKademliaID("0000000000000000000000000000000000000002")
	node1 := NewContact(kademliaID1, "")
	node2 := NewContact(kademliaID2, "")
	network := Network{}
	go network.Listen()
	go network.SendFindContactMessage(&node1, &node2)
	shortlis := <-network.shortlistCh
	if len(shortlis) != 0 {
		t.Errorf("Returned shortlist is not empty, got %d, want: %d", len(shortlis), 0)
	}
}
