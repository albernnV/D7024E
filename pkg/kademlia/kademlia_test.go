package kademlia

import (
	"testing"
)

var kademliaInstance *Kademlia = NewKademliaInstance(1, me)

func TestManageShortList(t *testing.T) {
	kademliaID1 := NewKademliaID("0000000000000000000000000000000000000001")
	kademliaID2 := NewKademliaID("0000000000000000000000000000000000000002")
	kademliaID3 := NewKademliaID("0000000000000000000000000000000000000003")
	contact1 := NewContact(kademliaID1, "")
	contact2 := NewContact(kademliaID2, "")
	contact3 := NewContact(kademliaID3, "")
	contact1.distance = NewKademliaID("0000000000000000000000000000000000000001")
	contact2.distance = NewKademliaID("0000000000000000000000000000000000000002")
	contact3.distance = NewKademliaID("0000000000000000000000000000000000000003")

	shortList := ContactCandidates{[]Contact{contact3, contact1, contact2, contact2}}
	newShortlist := []Contact{contact1}
	go func() {
		kademliaInstance.network.shortlistCh <- newShortlist //newshortlist gets written to achenneö
	}()
	kademliaInstance.manageShortlist(1, &shortList)

	cleanShortList := ContactCandidates{[]Contact{contact1, contact2, contact3}}

	if shortList.contacts[0] != cleanShortList.contacts[0] || shortList.contacts[1] != cleanShortList.contacts[1] || shortList.contacts[2] != cleanShortList.contacts[2] {
		t.Errorf("This shortlist was not managed, got: %v, want: %v", shortList.contacts, cleanShortList.contacts)
	}

}

func TestFindNotContacteNodes(t *testing.T) {
	kademliaID1 := NewKademliaID("0000000000000000000000000000000000000001")
	kademliaID2 := NewKademliaID("0000000000000000000000000000000000000002")
	contact1 := NewContact(kademliaID1, "")
	contact2 := NewContact(kademliaID2, "")
	contactedNodes := ContactCandidates{[]Contact{contact1}}
	shortList := ContactCandidates{[]Contact{contact1, contact2}}

	hasNotBeenContactedList := findNotContactedNodes(&shortList, &contactedNodes)

	hasNotBeenContactedListTest := ContactCandidates{[]Contact{contact2}}

	if hasNotBeenContactedList.contacts[0] != hasNotBeenContactedListTest.contacts[0] {
		t.Errorf("This does not find nodes that have not been contacted, got: %v, want: %v", hasNotBeenContactedList, hasNotBeenContactedListTest)
	}

}

func TestRemoveInactiveNodes(t *testing.T) {
	kademliaID1 := NewKademliaID("0000000000000000000000000000000000000001")
	kademliaID2 := NewKademliaID("0000000000000000000000000000000000000002")
	kademliaID3 := NewKademliaID("0000000000000000000000000000000000000003")
	contact1 := NewContact(kademliaID1, "")
	contact2 := NewContact(kademliaID2, "")
	contact3 := NewContact(kademliaID3, "")

	inactiveNodes := ContactCandidates{[]Contact{contact1, contact2}}
	shortList := ContactCandidates{[]Contact{contact1, contact2, contact3}}

	inactiveNodesList := removeInactiveNodes(shortList, inactiveNodes)

	cleanShortList := []Contact{contact3}

	if inactiveNodesList[0] != cleanShortList[0] {
		t.Errorf("This did not remove inactive nodes, got: %v, want: %v", inactiveNodesList, cleanShortList)
	}
}

func TestHashingData(t *testing.T) {

	data := []byte("Hejsan vad gör du?")

	hashedData := HashingData(data)

	sha1Hash := "c85373d0e75022b70dc94c99db4094ae80ab98d7" //sha1 hash for data

	if hashedData.String() != sha1Hash {
		t.Errorf("The hashed data is not correct, got: %s, want: %s", hashedData.String(), sha1Hash)
	}

}
