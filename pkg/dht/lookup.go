// pkg/dht/lookup.go
package dht

import (
	"sync"
)

const alpha = 3

// Lookup holds the state for a single iterative lookup process.
type Lookup struct {
	shortlist    *ContactCandidates
	queried      map[KademliaID]bool
	routingTable *RoutingTable
	rpc          RPC
	target       *KademliaID
}

// NewLookup creates a new Lookup instance.
func NewLookup(rt *RoutingTable, rpc RPC, target *KademliaID) *Lookup {
	return &Lookup{
		shortlist:    &ContactCandidates{},
		queried:      make(map[KademliaID]bool),
		routingTable: rt,
		rpc:          rpc,
		target:       target,
	}
}

// Start begins the iterative lookup process.
func (l *Lookup) Start() []Contact {
	// Start with the alpha closest nodes from our own routing table
	initialContacts := l.routingTable.FindClosestContacts(l.target, alpha)
	l.shortlist.Append(initialContacts)

	// Keep track of the closest contact found so far
	var closestContact *Contact
	if l.shortlist.Len() > 0 {
		// GetContacts returns a slice, so we take the first element.
		closestContact = &l.shortlist.GetContacts(1)[0]
	}

	// Main lookup loop
	for {
		contactsToQuery := l.getUnqueriedContacts(alpha)

		if len(contactsToQuery) == 0 {
			break
		}

		newContacts := l.queryContacts(contactsToQuery)
		l.shortlist.Append(newContacts)
		l.shortlist.Sort()

		if l.shortlist.Len() > 0 && (closestContact == nil || l.shortlist.GetContacts(1)[0].Less(closestContact)) {
			closestContact = &l.shortlist.GetContacts(1)[0]
		} else {
			// No closer contact found, so we are getting closer to the end.
			// Query the top k contacts that haven't been queried yet to be sure.
			remainingToQuery := l.getUnqueriedContacts(BucketSize)
			if len(remainingToQuery) > 0 {
				newContacts := l.queryContacts(remainingToQuery)
				l.shortlist.Append(newContacts)
				l.shortlist.Sort()
			}
			break
		}
	}

	return l.shortlist.GetContacts(BucketSize)
}

func (l *Lookup) getUnqueriedContacts(count int) []Contact {
	var contacts []Contact
	for _, contact := range l.shortlist.GetContacts(l.shortlist.Len()) {
		if len(contacts) >= count {
			break
		}
		if !l.queried[*contact.ID] {
			contacts = append(contacts, contact)
		}
	}
	return contacts
}

func (l *Lookup) queryContacts(contacts []Contact) []Contact {
	var newContacts []Contact
	var wg sync.WaitGroup
	resultsChan := make(chan []Contact, len(contacts))

	for _, contact := range contacts {
		l.queried[*contact.ID] = true
		wg.Add(1)
		go func(c Contact) {
			defer wg.Done()
			// Make sure the contact has its distance calculated relative to the target
			c.CalcDistance(l.target)
			
			foundContacts, err := l.rpc.FindNode(&c, l.target)
			if err != nil {
				return
			}
			resultsChan <- foundContacts
		}(contact)
	}

	wg.Wait()
	close(resultsChan)

	for contacts := range resultsChan {
		newContacts = append(newContacts, contacts...)
	}

	return newContacts
}
