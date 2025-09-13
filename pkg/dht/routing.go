// pkg/dht/routing.go
package dht

import (
	"container/heap"
	"time"

	"github.com/BrandonChongWenJun/D7024e-tutorial/pkg/id"
	"github.com/BrandonChongWenJun/D7024e-tutorial/pkg/proto"
)

// RoutingTable implements Kademlia K-buckets by distance range.
// We approximate the tree with 160 buckets, one per bit index of the MSB in XOR(self, other).
type RoutingTable struct {
	self    id.NodeID
	k       int
	buckets [160]*bucket
}

type bucket struct {
	entries []contact // LRU: index 0 is most recent
}

func newRoutingTable(self id.NodeID, k int) *RoutingTable {
	rt := &RoutingTable{self: self, k: k}
	for i := 0; i < len(rt.buckets); i++ {
		rt.buckets[i] = &bucket{entries: make([]contact, 0, k)}
	}
	return rt
}

// msbIndex returns the [0..159] index of the first set bit scanning from MSB to LSB.
// Returns -1 for zero distance (self).
func msbIndex(d [20]byte) int {
	for i := 0; i < 20; i++ {
		if d[i] == 0 { continue }
		b := d[i]
		for bit := 0; bit < 8; bit++ {
			if (b & (0x80 >> uint(bit))) != 0 {
				return i*8 + bit
			}
		}
	}
	return -1
}

func (rt *RoutingTable) bucketIndex(other id.NodeID) int {
	dist := id.XORDistance(rt.self, other)
	return msbIndex(dist)
}

// touch moves existing contact to front, or inserts if absent. If bucket overflows, drop LRU.
func (rt *RoutingTable) touch(pid id.NodeID, addr string) {
	idx := rt.bucketIndex(pid)
	if idx < 0 { return } // self
	b := rt.buckets[idx]
	now := time.Now()
	// find existing
	pos := -1
	for i, c := range b.entries {
		if c.addr == addr || c.id == pid {
			pos = i
			break
		}
	}
	if pos >= 0 {
		c := b.entries[pos]
		c.lastSeen = now
		c.addr = addr
		b.entries = append([]contact{c}, append(b.entries[:pos], b.entries[pos+1:]...)...)
		return
	}
	// insert new at front
	b.entries = append([]contact{{id: pid, addr: addr, lastSeen: now}}, b.entries...)
	if len(b.entries) > rt.k {
		// drop LRU (tail)
		b.entries = b.entries[:rt.k]
	}
}

// nearest returns up to k contacts closest to target.
func (rt *RoutingTable) nearest(target id.NodeID, k int) []proto.Contact {
	ph := &pairHeap{capK: k}
	heap.Init(ph)
	for _, b := range rt.buckets {
		for _, c := range b.entries {
			d := id.XORDistance(target, c.id)
			heap.Push(ph, pair{c: c, d: d})
			if ph.Len() > k {
				heap.Pop(ph)
			}
		}
	}
	// Drain heap to slice in ascending distance
	out := make([]proto.Contact, 0, k)
	tmp := make([]pair, 0, ph.Len())
	for ph.Len() > 0 { tmp = append(tmp, heap.Pop(ph).(pair)) }
	for i := len(tmp)-1; i >= 0; i-- {
		out = append(out, proto.Contact{ID: id.ToHex(tmp[i].c.id), Addr: tmp[i].c.addr})
	}
	return out
}

type pair struct { c contact; d [20]byte }
type pairHeap struct { items []pair; capK int }

func (h pairHeap) Len() int { return len(h.items) }
// Less implements max-heap (farther first) by returning true when i is farther than j.
func (h pairHeap) Less(i, j int) bool {
	di, dj := h.items[i].d, h.items[j].d
	for b := 0; b < 20; b++ {
		if di[b] != dj[b] { return di[b] > dj[b] }
	}
	return false
}
func (h pairHeap) Swap(i, j int) { h.items[i], h.items[j] = h.items[j], h.items[i] }
func (h *pairHeap) Push(x any) { h.items = append(h.items, x.(pair)) }
func (h *pairHeap) Pop() any {
	old := h.items
	n := len(old)
	it := old[n-1]
	h.items = old[:n-1]
	return it
}

// Node wrappers to use the routing table
func (n *Node) nearest(target id.NodeID, k int) []proto.Contact {
	return n.rt.nearest(target, k)
}

func (n *Node) insertContact(pid id.NodeID, addr string) {
	n.rt.touch(pid, addr)
}