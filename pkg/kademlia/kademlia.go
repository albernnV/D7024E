package kademlia

type Kademlia struct {
	alpha        int
	routingTable *RoutingTable
	network      *Network
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
	ch := make(chan []Contact)
	closerNodeHasBeenFound := true
	for closerNodeHasBeenFound {
		//Send find node RPC to alpha number of contacts in the shortlist
		for i := 0; i < kademlia.alpha; i++ {
			go SendFindNodeRPC(&shortlist.contacts[i], kademlia.network, ch)
		}

		//Retrieve all shortlists from the contacted nodes
		subShortlist1, subShortlist2, subShortlist3 := <-ch, <-ch, <-ch
		s := append(subShortlist1, subShortlist2...)
		newShortList := ContactCandidates{append(s, subShortlist3...)}
		//Sort the new shortlist and remove any duplicates
		newShortList.Sort()
		// TODO: Remove dublicates
		shortlist = newShortList
		if shortlist.contacts[0].Less(closestContact) {
			closestContact = &shortlist.contacts[0]
		} else {
			closerNodeHasBeenFound = false
			//Send a RPC to each of the k closest nodes that has not already been queried
			//Stop FIND_NODE_RPC
		}
	}
	return &ContactCandidates{}
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
func SendFindNodeRPC(contact *Contact, network *Network, channel chan []Contact) {
	closestNodes := network.SendFindContactMessage(contact)
	channel <- closestNodes
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
	// TODO
}
