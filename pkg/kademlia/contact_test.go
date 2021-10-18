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
	dist := NewKademliaID("0000000000000000000000000000000000000005")

	newContact := NewContact(kademliaID, addr)
	newContact.distance = dist
	contactToString := newContact.String()

	testString := `contact("0000000000000000000000000000000000000001", "172.0.0.1:8000", "0000000000000000000000000000000000000005")`

	if contactToString != testString {
		t.Errorf("This does not output the same string, got: %s, want: %s", contactToString, testString)
	}

	testString = `contact("0000000000000000000000000000000000000001", "172.0.0.1:8000", "")`
	contactNilDistance := NewContact(kademliaID, addr)
	contactToString = contactNilDistance.String()

	if contactToString != testString {
		t.Errorf("This does not output the same string, got: %s, want: %s", contactToString, testString)
	}
}

func TestAppend(t *testing.T) {

	kademliaID := NewKademliaID("0000000000000000000000000000000000000001")
	testKademliaID := NewKademliaID("0000000000000000000000000000000000000001")

	contact := NewContact(kademliaID, "172.0.0.1:8000")
	testContact := NewContact(testKademliaID, "172.0.0.1:8000")

	contactCandidates := ContactCandidates{}
	contactCandidates.Append([]Contact{contact})

	testContactCandidates := ContactCandidates{[]Contact{testContact}}

	if *contactCandidates.contacts[0].ID != *testContactCandidates.contacts[0].ID || contactCandidates.contacts[0].Address != testContactCandidates.contacts[0].Address {
		t.Errorf("This does not append to the ContactCandidates, got: %v and %v, want: %v and %v",
			contactCandidates.contacts[0].ID, contactCandidates.contacts[0].Address, testContactCandidates.contacts[0].ID, testContactCandidates.contacts[0].Address)
	}

}

func TestGetContacts(t *testing.T) {

	kademliaID1 := NewKademliaID("0000000000000000000000000000000000000001")
	newContact := NewContact(kademliaID1, "172.0.0.1:8000")
	contactCandidates := ContactCandidates{[]Contact{newContact}}

	getContact1 := contactCandidates.GetContacts(1)
	getContact2 := contactCandidates.GetContacts(2)

	if getContact1[0] != newContact || getContact2[0] != newContact {
		t.Errorf("This does not return the correct contact, got: %v and %v, want: %v", getContact1[0], getContact2[0], newContact)
	}

}

func TestSort(t *testing.T) {

	kademliaID1 := NewKademliaID("0000000000000000000000000000000000000001")
	contact1 := NewContact(kademliaID1, "172.0.0.1:8000")
	kademliaID2 := NewKademliaID("0000000000000000000000000000000000000002")
	contact2 := NewContact(kademliaID2, "172.0.0.2:8000")
	kademliaID3 := NewKademliaID("0000000000000000000000000000000000000003")
	contact3 := NewContact(kademliaID3, "172.0.0.3:8000")

	contact1.distance = NewKademliaID("0000000000000000000000000000000000000001")
	contact2.distance = NewKademliaID("0000000000000000000000000000000000000002")
	contact3.distance = NewKademliaID("0000000000000000000000000000000000000003")

	contactCandidates := ContactCandidates{[]Contact{contact3, contact1, contact2}}
	contactCandidates.Sort()

	sortedContactCandidates := ContactCandidates{[]Contact{contact1, contact2, contact3}}

	if contactCandidates.contacts[0] != sortedContactCandidates.contacts[0] || contactCandidates.contacts[1] != sortedContactCandidates.contacts[1] || contactCandidates.contacts[2] != sortedContactCandidates.contacts[2] {
		t.Errorf("This condactCandidates is not sorted, got: %s, want: %s", contactCandidates.contacts[0].String(), sortedContactCandidates.contacts[0].String())
	}

}

func TestContactLen(t *testing.T) {

	kademliaID1 := NewKademliaID("0000000000000000000000000000000000000001")
	contact1 := NewContact(kademliaID1, "172.0.0.1:8000")
	kademliaID2 := NewKademliaID("0000000000000000000000000000000000000002")
	contact2 := NewContact(kademliaID2, "172.0.0.2:8000")
	kademliaID3 := NewKademliaID("0000000000000000000000000000000000000003")
	contact3 := NewContact(kademliaID3, "172.0.0.3:8000")

	contactCandidates := ContactCandidates{[]Contact{contact1, contact2, contact3}}

	if contactCandidates.Len() != 3 {
		t.Errorf("This is not the right length, got: %v, want: %v", contactCandidates, 3)
	}
}

func TestSwap(t *testing.T) {

	kademliaID1 := NewKademliaID("0000000000000000000000000000000000000001")
	contact1 := NewContact(kademliaID1, "172.0.0.1:8000")
	kademliaID2 := NewKademliaID("0000000000000000000000000000000000000002")
	contact2 := NewContact(kademliaID2, "172.0.0.2:8000")
	kademliaID3 := NewKademliaID("0000000000000000000000000000000000000003")
	contact3 := NewContact(kademliaID3, "172.0.0.3:8000")

	contactCandidates := ContactCandidates{[]Contact{contact1, contact2, contact3}}
	contactCandidates.Swap(0, 1)

	swappedContactCandidates := ContactCandidates{[]Contact{contact2, contact1, contact3}}

	if contactCandidates.contacts[0].ID != swappedContactCandidates.contacts[0].ID || contactCandidates.contacts[1].ID != swappedContactCandidates.contacts[1].ID {
		t.Errorf("This did not swap the contacts, got: %v and %v, want: %v and %v", contactCandidates.contacts[0], contactCandidates.contacts[1], swappedContactCandidates.contacts[0], swappedContactCandidates.contacts[1])
	}

}

func TestContactCandidatesLess(t *testing.T) {
	kademliaID1 := NewKademliaID("0000000000000000000000000000000000000001")
	contact1 := NewContact(kademliaID1, "172.0.0.1:8000")
	kademliaID2 := NewKademliaID("0000000000000000000000000000000000000002")
	contact2 := NewContact(kademliaID2, "172.0.0.2:8000")
	kademliaID3 := NewKademliaID("0000000000000000000000000000000000000003")
	contact3 := NewContact(kademliaID3, "172.0.0.3:8000")

	contact1.distance = NewKademliaID("0000000000000000000000000000000000000001")
	contact2.distance = NewKademliaID("0000000000000000000000000000000000000002")
	contact3.distance = NewKademliaID("0000000000000000000000000000000000000003")

	contactCandidates := ContactCandidates{[]Contact{contact1, contact2, contact3}}

	if contactCandidates.Less(0, 1) != true {
		t.Errorf("This is not smaller than the preceding index, got: %v", contactCandidates.contacts)
	}

}

func TestRemoveDuplicates(t *testing.T) {

	kademliaID1 := NewKademliaID("0000000000000000000000000000000000000001")
	contact1 := NewContact(kademliaID1, "172.0.0.1:8000")
	kademliaID2 := NewKademliaID("0000000000000000000000000000000000000001")
	contact2 := NewContact(kademliaID2, "172.0.0.2:8000")
	kademliaID3 := NewKademliaID("0000000000000000000000000000000000000002")
	contact3 := NewContact(kademliaID3, "172.0.0.3:8000")

	contactCandidates := ContactCandidates{[]Contact{contact1, contact2, contact3}}
	contactCandidates.RemoveDuplicates()

	contactCandidatesNoDuplicates := ContactCandidates{[]Contact{contact1, contact3}}

	if contactCandidates.contacts[0].ID != contactCandidatesNoDuplicates.contacts[0].ID || contactCandidates.contacts[1].ID != contactCandidatesNoDuplicates.contacts[1].ID {
		t.Errorf("This does not remove duplicates, got: %v, want: %v", contactCandidates.contacts, contactCandidatesNoDuplicates.contacts)
	}

}
