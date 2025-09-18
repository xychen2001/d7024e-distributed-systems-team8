// pkg/dht/bucket.go
package dht

import (
	"container/list"
)

// bucket definition
// contains a List
type bucket struct {
	list *list.List
}

// newBucket returns a new instance of a bucket
func newBucket() *bucket {
	bucket := &bucket{}
	bucket.list = list.New()
	return bucket
}

// AddContact adds a new contact to the bucket.
// It follows the LRU discipline: if the contact already exists, it's moved to the front.
// If the bucket is full, the new contact is not added.
func (bucket *bucket) AddContact(contact Contact) {
	var element *list.Element
	// Find the contact in the bucket
	for e := bucket.list.Front(); e != nil; e = e.Next() {
		if (contact).ID.Equals(e.Value.(Contact).ID) {
			element = e
		}
	}

	if element != nil {
		// If the contact already exists, move it to the front (most recently seen).
		bucket.list.MoveToFront(element)
	} else {
		// If the contact does not exist, add it to the front if there is space.
		if bucket.list.Len() < bucketSize {
			bucket.list.PushFront(contact)
		} else {
			// TODO: Implement the eviction policy from the sprint plan.
			// (Ping the least-recently-seen contact and evict on timeout).
		}
	}
}

// GetContactAndCalcDistance returns an array of Contacts where
// the distance has already been calculated
func (bucket *bucket) GetContactAndCalcDistance(target *KademliaID) []Contact {
	var contacts []Contact

	for elt := bucket.list.Front(); elt != nil; elt = elt.Next() {
		contact := elt.Value.(Contact)
		contact.CalcDistance(target)
		contacts = append(contacts, contact)
	}

	return contacts
}

// Len return the size of the bucket
func (bucket *bucket) Len() int {
	return bucket.list.Len()
}


