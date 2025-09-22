package dht

import (
	"errors"
	"testing"
)

// mockRPC is a mock implementation of the RPC interface for testing.
type mockRPC struct {
	pingShouldFail bool
}

func (m *mockRPC) FindNode(contact *Contact, target *KademliaID) ([]Contact, error) {
	// Not needed for this test
	return nil, nil
}

func (m *mockRPC) Ping(contact *Contact) error {
	if m.pingShouldFail {
		return errors.New("ping failed")
	}
	return nil
}

func TestBucketEviction(t *testing.T) {
	// Test Case 1: Ping fails, least recently seen contact should be evicted.
	t.Run("Ping Fails", func(t *testing.T) {
		bucket := newBucket()
		mockRPC := &mockRPC{pingShouldFail: true}

		// Fill the bucket to its capacity.
		for i := 0; i < bucketSize; i++ {
			contact := NewContact(NewRandomKademliaID(), "")
			bucket.AddContact(contact, mockRPC)
		}

		lruContact := bucket.list.Back().Value.(Contact)
		newContact := NewContact(NewRandomKademliaID(), "")

		bucket.AddContact(newContact, mockRPC)

		// Check if the new contact was added.
		found := false
		for e := bucket.list.Front(); e != nil; e = e.Next() {
			if e.Value.(Contact).ID.Equals(newContact.ID) {
				found = true
				break
			}
		}
		if !found {
			t.Error("New contact should have been added to the bucket")
		}

		// Check if the least recently seen contact was evicted.
		found = false
		for e := bucket.list.Front(); e != nil; e = e.Next() {
			if e.Value.(Contact).ID.Equals(lruContact.ID) {
				found = true
				break
			}
		}
		if found {
			t.Error("Least recently seen contact should have been evicted")
		}
	})

	// Test Case 2: Ping succeeds, least recently seen contact should be moved to the front.
	t.Run("Ping Succeeds", func(t *testing.T) {
		bucket := newBucket()
		mockRPC := &mockRPC{pingShouldFail: false}

		// Fill the bucket to its capacity.
		for i := 0; i < bucketSize; i++ {
			contact := NewContact(NewRandomKademliaID(), "")
			bucket.AddContact(contact, mockRPC)
		}

		lruContact := bucket.list.Back().Value.(Contact)
		newContact := NewContact(NewRandomKademliaID(), "")

		bucket.AddContact(newContact, mockRPC)

		// Check if the new contact was not added.
		found := false
		for e := bucket.list.Front(); e != nil; e = e.Next() {
			if e.Value.(Contact).ID.Equals(newContact.ID) {
				found = true
				break
			}
		}
		if found {
			t.Error("New contact should not have been added to the bucket")
		}

		// Check if the least recently seen contact was moved to the front.
		if !bucket.list.Front().Value.(Contact).ID.Equals(lruContact.ID) {
			t.Error("Least recently seen contact should have been moved to the front")
		}
	})
}
