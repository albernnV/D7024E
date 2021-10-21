package kademlia

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
)

type Kademlia struct {
	alpha   int
	network *Network
}

//Starting function for the nodes by pinging the bootstrap node to be able to join the network
func (kademlia *Kademlia) Start() {
	// Join network by perforing lookup on yourself
	bootstrapNode := NewContact(nil, "172.18.0.3:8000")
	if kademlia.network.routingTable.me.Address != bootstrapNode.Address {
		kademlia.network.SendPingMessage(&bootstrapNode)
	}
	go kademlia.network.LoopListen()
	kademlia.LookupContact(&kademlia.network.routingTable.me)
}

// Returns a kademlia instance to be able to use it as a reference
func NewKademliaInstance(alpha int, me Contact) *Kademlia {
	network := NewNetwork(me)
	newKademliaInstance := &Kademlia{alpha, network}
	return newKademliaInstance
}

//Returns k closests nodes in a shortlist which includes cleaning up the shortlist and sending find contact message to the nodes
func (kademlia *Kademlia) LookupContact(target *Contact) *ContactCandidates {
	//Find k closest nodes
	closestNodes := kademlia.network.routingTable.FindClosestContacts(target.ID, kademlia.alpha)
	if len(closestNodes) != 0 {
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
			contactsToSendTo := shortlist.GetContacts(kademlia.alpha)
			for _, contact := range contactsToSendTo {
				go kademlia.network.SendFindContactMessage(&contact, target)
				hasBeenContactedList.Append([]Contact{contact})
			}
			kademlia.manageShortlist(len(contactsToSendTo), &shortlist)
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
				kademlia.manageShortlist(nodesToContact.Len(), &shortlist)
				//Remove all inactive nodes from the shortlist
				shortlist.contacts = removeInactiveNodes(shortlist, kademlia.network.inactiveNodes)
			}
		}

		return &shortlist
	} else {
		return &ContactCandidates{[]Contact{}}
	}
}

//Manages the shortlist by sorting, removing the duplicates, removes "itself" and keeps k size of the shortlist
func (kademlia *Kademlia) manageShortlist(alpha int, shortlist *ContactCandidates) {
	for i := 0; i < alpha; i++ {
		newShortList := <-kademlia.network.shortlistCh
		shortlist.Append(newShortList)
		shortlist.Sort()
		shortlist.RemoveDuplicates()
		shortlist.RemoveContact(&kademlia.network.routingTable.me) //Remove self from shortlist
		shortlist.contacts = shortlist.GetContacts(bucketSize)
	}
}

func (kademlia *Kademlia) retreiveValues(alpha int, valueSoughtAfter string) ValueAndSender {
	var currentValue ValueAndSender
	for i := 0; i < alpha; i++ {
		receivedData := <-kademlia.network.receviedValueChan
		if receivedData.value != currentValue.value && receivedData.value == valueSoughtAfter {
			currentValue = receivedData
		}
	}
	return currentValue
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

//Returns a clean shortlist without the inactive nodes
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

//LookupData fetches the hashed data object from the closests nodes and sends a find data message
func (kademlia *Kademlia) LookupData(data string) {
	// Make the hash into a kademliaID to be able to make a new contact
	hashToKademliaID := HashingData([]byte(data))
	kademliaToContact := NewContact(hashToKademliaID, "")
	// Look up the closests contacts
	shortlist := kademlia.LookupContact(&kademliaToContact)

	//loop through all contact and find value
	for _, nodeToContact := range shortlist.contacts {
		go kademlia.network.SendFindDataMessage(hashToKademliaID.String(), nodeToContact)
	}
	receivedData := kademlia.retreiveValues(shortlist.Len(), data)
	fmt.Println("Received " + receivedData.value + " from: " + receivedData.sender.ID.String())
}

//Store finds the closests nodes for the hashed data object and sends a store message to those nodes
func (kademlia *Kademlia) Store(data []byte) {
	//Hash the data to get a newKademliaID
	fileKademliaID := HashingData(data)
	newContact := NewContact(fileKademliaID, "")

	// Find closest contacts for the key
	closestsNodes := kademlia.LookupContact(&newContact)
	//SendStore RPCs
	for _, nodeToStoreAt := range closestsNodes.contacts {
		kademlia.network.SendStoreMessage(data, &nodeToStoreAt)
	}
	fmt.Println("hash: " + fileKademliaID.String())
}

//Hashes the data and returns the hash as a NewKademliaID
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
