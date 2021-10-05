package kademlia

import (
	"crypto/sha1"
	"encoding/hex"
)

type Kademlia struct {
	alpha        int
	routingTable *RoutingTable
	network      *Network
}

func (kademlia *Kademlia) Start() {
	go kademlia.routingTable.UpdateRoutingTable()
}

func (kademlia *Kademlia) Stop() {
	close(kademlia.routingTable.routingTableChan)
}

func NewKademliaInstance(alpha int, me Contact) *Kademlia {
	routingTable := NewRoutingTable(me)
	network := &Network{}
	newKademliaInstance := &Kademlia{alpha, routingTable, network}
	return newKademliaInstance
}

func (kademlia *Kademlia) LookupContact(target *Contact) *ContactCandidates {
	// TODO
	//FindClosestContacts() for target
	//
	//Check which node is closest to target, save in closestNode
	//
	//Input all alpha=3 contacts into a shortlist
	//
	//SendFindContactMessage() in parallel to the alpha nodes in the shortlist. These chould return k contacts. If anyone fails
	//to reply they are removed from the shortlist
	//
	//The node then fills the shortlist with contacts from the replies received. These are those closest to the target.
	//
	//update closestNode
	//
	//SendFindContactMessage() in parallel to another alpha nodes from the shortlist. The condition is that they haven't been
	//contacted already.
	//
	//continue until either no node in the returned sets are closer to the closest node seen or the initiating node has
	//accumulated k probed and known to be active contacts.
	//
	//If a cycle doesn't find a node that is closer than the already closest node the node will send a RPC each of
	//the k closest nodes that it has not already queried

	//Find k closest nodes
	closestNodes := kademlia.routingTable.FindClosestContacts(target.ID, kademlia.alpha)
	//Initiate closestNode
	closestContact := findClosestContact(closestNodes, target)
	//Initiate shortlist
	var shortlist ContactCandidates
	shortlist.contacts = closestNodes
	//Initiate channel where shortlists from the goroutines will be written to
	shortlistCh := make(chan []Contact)
	//Channel for writing inactive nodes to
	hasNotAnsweredCh := make(chan Contact)

	hasBeenContactedList := ContactCandidates{}
	hasNotAnsweredList := ContactCandidates{}

	closerNodeHasBeenFound := true
	go manageInactiveNodes(hasNotAnsweredCh, &hasNotAnsweredList)
	for closerNodeHasBeenFound {
		//Send find node RPC to alpha number of contacts in the shortlist
		for i := 0; i < kademlia.alpha; i++ {
			go kademlia.SendFindNodeRPC(&shortlist.contacts[i], target, kademlia.network, shortlistCh, hasNotAnsweredCh)
		}
		manageShortlist(kademlia.alpha, &shortlist, shortlistCh)
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
			for _, nodeToContact := range nodesToContact.GetContacts(bucketSize) {
				go kademlia.SendFindNodeRPC(&nodeToContact, target, kademlia.network, shortlistCh, hasNotAnsweredCh)
			}
			manageShortlist(kademlia.alpha, &shortlist, shortlistCh)
			//Remove all inactive nodes from the shortlist
			shortlist.contacts = removeInactiveNodes(shortlist, hasNotAnsweredList)
		}
	}

	return &shortlist
}

//Sends a find node RPC to the contact which will send back the k closest nodes. These contacts will be written to the
//channel to be retireved
func (kademlia *Kademlia) SendFindNodeRPC(contact *Contact, target *Contact, network *Network, shortlistChannel chan []Contact, hasNotAnsweredChannel chan Contact) {
	closestNodes, didNotAnswer := network.SendFindContactMessage(contact, target, hasNotAnsweredChannel)
	shortlistChannel <- closestNodes
	if didNotAnswer {
		hasNotAnsweredChannel <- *contact
	} else {
		//Add contact to routing table
		kademlia.routingTable.routingTableChan <- *contact
	}
}

func manageShortlist(alpha int, shortlist *ContactCandidates, shortlistCh chan []Contact) {
	for i := 0; i < alpha; i++ {
		newShortList := <-shortlistCh
		shortlist.Append(newShortList)
		shortlist.Sort()
		shortlist.RemoveDuplicates()
		shortlist.contacts = shortlist.GetContacts(bucketSize)
	}

}

//Calculates the distances from the contacts to the target contact and returns the contact with the shortest distance
func findClosestContact(contacts []Contact, target *Contact) *Contact {
	var closestNode *Contact = &contacts[0]
	for i := 0; i < len(contacts); i++ {
		contacts[i].CalcDistance(target.ID)
	}
	//Compare distances with the closestNode and update it accordingly
	for j := 0; j < len(contacts); j++ {
		if contacts[j].distance.Less(closestNode.distance) {
			closestNode = &contacts[j]
		}
	}

	return closestNode
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

func manageInactiveNodes(hasNotAnsweredCh chan Contact, hasNotAnsweredList *ContactCandidates) {
	for {
		inactiveNode := <-hasNotAnsweredCh
		hasNotAnsweredList.Append([]Contact{inactiveNode})
	}
}

//Removes all inactive nodes in inactiveNodes from shortlist
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
	// TODO
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
	stringToBytes.Write(data)
	hashedData := stringToBytes.Sum(nil)
	// Encodes the hash back to string to make it a new kademlia ID
	hashedStringData := hex.EncodeToString(hashedData)
	hashedKademliaID := NewKademliaID(hashedStringData)

	return hashedKademliaID
}
