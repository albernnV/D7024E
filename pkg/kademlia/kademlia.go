package kademlia

import (
	"crypto/sha1"
	"encoding/hex"
)

type Kademlia struct {
	alpha   int
	network *Network
}

func (kademlia *Kademlia) Start() {
	go kademlia.network.routingTable.UpdateRoutingTable()
}

func (kademlia *Kademlia) Stop() {
	close(kademlia.network.routingTable.routingTableChan)
}

func NewKademliaInstance(alpha int, me Contact) *Kademlia {
	network := NewNetwork(me)
	newKademliaInstance := &Kademlia{alpha, network}
	return newKademliaInstance
}

func (kademlia *Kademlia) LookupContact(target *Contact) *ContactCandidates {
	//Find k closest nodes
	closestNodes := kademlia.network.routingTable.FindClosestContacts(target.ID, kademlia.alpha)
	//Initiate closestNode
	closestNodesToContactCandidates := ContactCandidates{closestNodes}
	closestNodesToContactCandidates.Sort()
	closestContact := &closestNodesToContactCandidates.contacts[0]
	//Initiate shortlist
	var shortlist ContactCandidates
	shortlist.contacts = closestNodes

	hasBeenContactedList := ContactCandidates{}

	closerNodeHasBeenFound := true
	for closerNodeHasBeenFound {
		//Send find node RPC to alpha number of contacts in the shortlist
		for i := 0; i < kademlia.alpha; i++ {
			go kademlia.network.SendFindContactMessage(&shortlist.contacts[i], target)
		}
		kademlia.manageShortlist(kademlia.alpha, &shortlist)
		//Check end condition
		if shortlist.contacts[0].Less(closestContact) {
			closestContact = &shortlist.contacts[0]
		} else {
			closerNodeHasBeenFound = false
			//Find closest nodes that have not yet been contacted
			nodesToContact := findNotContactedNodes(&shortlist, &hasBeenContactedList)
			nodesToContact.Sort()
			nodesToContact.RemoveDuplicates()
			//Send a RPC to each of the k closest nodes that has not already been contacted
			for _, nodeToContact := range nodesToContact.contacts {
				go kademlia.network.SendFindContactMessage(&nodeToContact, target)
			}
			kademlia.manageShortlist(bucketSize, &shortlist)
			//Remove all inactive nodes from the shortlist
			shortlist.contacts = removeInactiveNodes(shortlist, kademlia.network.inactiveNodes)
		}
	}

	return &shortlist
}

func (kademlia *Kademlia) manageShortlist(alpha int, shortlist *ContactCandidates) {
	for i := 0; i < alpha; i++ {
		newShortList := <-kademlia.network.shortlistCh
		shortlist.Append(newShortList)
		shortlist.Sort()
		shortlist.RemoveDuplicates()
		shortlist.contacts = shortlist.GetContacts(bucketSize)
	}

}

//Returns the contacts in the shortlis that haven't been contacted
func findNotContactedNodes(shortlist *ContactCandidates, contactedNodes *ContactCandidates) ContactCandidates {
	hasNotBeenContactedList := make([]Contact, 0)
	for _, contact := range shortlist.contacts {
		hasNotBeenContacted := true
		for _, contactedNode := range contactedNodes.contacts {
			if contact.ID == contactedNode.ID {
				hasNotBeenContacted = false
			}
		}
		if hasNotBeenContacted {
			hasNotBeenContactedList = append(hasNotBeenContactedList, contact)
		}
	}
	return ContactCandidates{hasNotBeenContactedList}
}

func removeInactiveNodes(shortlist ContactCandidates, inactiveNodes ContactCandidates) []Contact {
	cleanShortlist := make([]Contact, 0)
	for _, contact := range shortlist.contacts {
		isActive := true
		for _, inactiveNode := range inactiveNodes.contacts {
			if contact.ID == inactiveNode.ID {
				isActive = false
			}
		}
		if isActive {
			cleanShortlist = append(cleanShortlist, contact)
		}
	}
	return cleanShortlist
}

func (kademlia *Kademlia) LookupData(hash string) {
	// Make the hash into a kademliaID to be able to make a new contact
	hashToKademliaID := NewKademliaID(hash)
	kademliaToContact := NewContact(hashToKademliaID, "")
	// Look up the closests contacts
	shortlist := kademlia.LookupContact(&kademliaToContact)

	//loop through all contact and find value
	for _, nodeToContact := range shortlist.contacts {
		kademlia.network.SendFindDataMessage(hash, &nodeToContact)
	}
}

func (kademlia *Kademlia) Store(data []byte) {
	//Hash the data to get a newKademliaID
	fileKademliaID := HashingData(data)
	newContact := NewContact(fileKademliaID, "")

	// Find closest contacts for the key
	closestsNodes := kademlia.LookupContact(&newContact)

	//SendStore RPCs
	for _, nodeToStoreAt := range closestsNodes.contacts {
		go kademlia.network.SendStoreMessage(data, &nodeToStoreAt, newContact)
	}

}

func HashingData(data []byte) *KademliaID {
	//hash the data
	stringToBytes := sha1.New()
	stringToBytes.Write([]byte(data))
	hashedData := stringToBytes.Sum(nil)

	// Encodes the hash back to string to make it a new kademlia ID
	hashedStringData := hex.EncodeToString(hashedData)
	hashedKademliaID := NewKademliaID(hashedStringData)

	return hashedKademliaID
}
