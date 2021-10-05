package kademlia

import (
	"testing"
)

/*func TestNewKademliaID(t *testing.T) {

	s := "00000000000000000001"
	byteString := []byte(s)

	newKademliaID := NewKademliaID(s)

	if newKademliaID != byteString {
		t.Errorf("This is not a KademliaID, got: %v, want: %v", newKademliaID, byteString)
	}
}*/

func TestString(t *testing.T) {

	s := "0000000000000000000000000000000000000001"
	newKademliaID := NewKademliaID(s)

	if newKademliaID.String() != s {
		t.Errorf("This does not return a simple string represetation of a KademliaID, got: %s, want: %s", newKademliaID.String(), s)
	}
}

func TestCalcDistance(t *testing.T) {

	ID1 := "0000000000000000000000000000000000000001"
	ID2 := "0000000000000000000000000000000000000002"
	newKademliaID1 := NewKademliaID(ID1)
	newKademliaID2 := NewKademliaID(ID2)
	newKademliaID := newKademliaID1.CalcDistance(newKademliaID2)

	ID1xorID2 := NewKademliaID("0000000000000000000000000000000000000003")

	if *newKademliaID != *ID1xorID2 {
		t.Errorf("This does not calculate the distance between 2 KademliaIDs, got: %v, want: %v", *newKademliaID, *ID1xorID2)
	}
}
