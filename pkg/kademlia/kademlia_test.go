package kademlia

import (
	"testing"
)

func TestSendFindNodeRPC(t *testing.T) {
	kademliaID1 := NewKademliaID("0000000000000000000000000000000000000001")
	node1 := NewContact(kademliaID1, "")

	kademliaInstance := NewKademliaInstance(3, node1)

	kademliaID2 := NewKademliaID("0000000000000000000000000000000000000002")
	node2 := NewContact(kademliaID2, "")
	go kademliaInstance.network.Listen()
	go kademliaInstance.network.SendFindContactMessage(&node1, &node2)
	shortlis := <-kademliaInstance.network.shortlistCh
	if len(shortlis) != 0 {
		t.Errorf("Returned shortlist is not empty, got: %d, want: %d", len(shortlis), 0)
	}
}

func TestHashingData(t *testing.T) {

	s := []byte("Hejsan vad gör du?")

	hashedData := HashingData(s)

	sha1Hash := "c85373d0e75022b70dc94c99db4094ae80ab98d7"

	if hashedData.String() != sha1Hash {
		t.Errorf("The hashed data is not correct, got: %s, want: %s", hashedData.String(), sha1Hash)
	}

}
