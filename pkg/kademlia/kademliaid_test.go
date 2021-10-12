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
func TestNewRandomKademliaID(t *testing.T) {

	randomKademliaID1 := NewRandomKademliaID()
	randomKademliaID2 := NewRandomKademliaID()
	randomKademliaID3 := NewRandomKademliaID()

	randomKademliaIDs := []*KademliaID{randomKademliaID1, randomKademliaID2, randomKademliaID3}
	newRandomKademliaIDs := make([]*KademliaID, 0)

	for _, idToAdd := range randomKademliaIDs {
		for _, id := range newRandomKademliaIDs {
			if *id == *idToAdd {
				t.Errorf("This does not return a random KademliaID")
			}
		}
		newRandomKademliaIDs = append(newRandomKademliaIDs, idToAdd)
	}

}

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

func TestEquals(t *testing.T) {
	newKademliaID1 := NewKademliaID("0000000000000000000000000000000000000001")
	testKademliaID := NewKademliaID("0000000000000000000000000000000000000001")

	newKademliaID := newKademliaID1.Equals(testKademliaID)

	if newKademliaID != true {
		t.Errorf("This is not equals the testKademliaID, got %v, want: %v", newKademliaID, testKademliaID)
	}
}

func TestLess(t *testing.T) {

	kademliaID := NewKademliaID("0000000000000000000000000000000000000001")
	otherKademliaID := NewKademliaID("0000000000000000000000000000000000000005")

	newKademliaID1 := kademliaID.Less(otherKademliaID)
	newKademliaID2 := otherKademliaID.Less(kademliaID)

	if newKademliaID1 != true || newKademliaID2 != false {
		t.Errorf("This Less() does not work, got: %v and %v, want: %v", newKademliaID1, newKademliaID2, otherKademliaID)
	}

}
