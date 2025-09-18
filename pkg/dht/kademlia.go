// pkg/dht/kademlia.go
package dht

// Kademlia represents a Kademlia node.
type Kademlia struct {
	RoutingTable *RoutingTable
	Network      RPC
}

// NewKademlia creates a new Kademlia instance.
func NewKademlia(rt *RoutingTable, rpc RPC) *Kademlia {
	return &Kademlia{
		RoutingTable: rt,
		Network:      rpc,
	}
}

// LookupContact performs the iterative lookup process to find the k closest contacts to the target.
func (k *Kademlia) LookupContact(target *KademliaID) []Contact {
	lookup := NewLookup(k.RoutingTable, k.Network, target)
	return lookup.Start()
}
