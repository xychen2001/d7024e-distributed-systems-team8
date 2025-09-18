// pkg/dht/net.go
package dht

// RPC is an interface for making network requests to other Kademlia nodes.
type RPC interface {
	// FindNode sends a FIND_NODE request to a contact and returns a list of closer contacts.
	FindNode(contact *Contact, target *KademliaID) ([]Contact, error)
	// Ping sends a PING request to a contact and expects a PONG in return.
	Ping(contact *Contact) error
}
