package kademlia

type Kademlia struct {
	alpha        int
	routingTable *RoutingTable
	network      *Network
}

func (kademlia *Kademlia) LookupContact(target *Contact) {
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
	hasNotAnsweredCh := make(chan Contact)

	hasBeenContactedList := ContactCandidates{}
	hasNotAnsweredList := ContactCandidates{}

	closerNodeHasBeenFound := true
	for closerNodeHasBeenFound {
		//Send find node RPC to alpha number of contacts in the shortlist
		for i := 0; i < kademlia.alpha; i++ {
			go SendFindNodeRPC(&shortlist.contacts[i], target, kademlia.network, shortlistCh, hasNotAnsweredCh)
		}
		//Add contacted nodes to the list
		hasBeenContactedList.Append(shortlist.GetContacts(kademlia.alpha))
		//Retrieve all contacts that failed to answered
		inactiveNodes := retrieveInactiveNodes(kademlia.alpha, hasNotAnsweredCh)
		hasNotAnsweredList.Append(inactiveNodes)
		//Retrieve all shortlists from the contacted nodes
		for i := 0; i < kademlia.alpha; i++ {
			newShortList := <-shortlistCh
			shortlist.Append(newShortList)
		}
		//Sort the new shortlist and remove any duplicates
		shortlist.Sort()
		shortlist.RemoveDuplicates()
		shortlist.contacts = shortlist.GetContacts(bucketSize)
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
				go SendFindNodeRPC(&nodeToContact, target, kademlia.network, shortlistCh, hasNotAnsweredCh)
			}
			//Retrieve shortlists from contacted nodes
			for i := 0; i < bucketSize; i++ {
				newShortList := <-shortlistCh
				shortlist.Append(newShortList)
			}
			shortlist.Sort()
			shortlist.RemoveDuplicates()
			//Remove inactive nodes

			shortlist.contacts = shortlist.GetContacts(bucketSize)
			//Stop FIND_NODE_RPC
		}
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

//Sends a find node RPC to the contact which will send back the k closest nodes. These contacts will be written to the
//channel to be retireved
func SendFindNodeRPC(contact *Contact, target *Contact, network *Network, shortlistChannel chan []Contact, hasNotAnsweredChannel chan Contact) {
	closestNodes := network.SendFindContactMessage(contact, target, hasNotAnsweredChannel)
	shortlistChannel <- closestNodes
}

func retrieveInactiveNodes(alpha int, hasNotAnsweredCh chan Contact) []Contact {
	inactiveNodes := make([]Contact, bucketSize)
	for i := 0; i < alpha; i++ {
		inactiveNode := <-hasNotAnsweredCh
		inactiveNodes = append(inactiveNodes, inactiveNode)
	}
	return inactiveNodes
}

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

func (kademlia *Kademlia) LookupData(hash string) {
	// TODO
}

func (kademlia *Kademlia) Store(data []byte) {
	// TODO
}
