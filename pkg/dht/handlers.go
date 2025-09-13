// pkg/dht/handlers.go
package dht

import (
	"net"
	"time"

	"github.com/BrandonChongWenJun/D7024e-tutorial/pkg/id"
	"github.com/BrandonChongWenJun/D7024e-tutorial/pkg/proto"
)

func (n *Node) handle(m proto.Message, src *net.UDPAddr) {
	switch m.Type {
	case proto.MsgPing:
		n.touch(m.From.ID, m.From.Addr)
		n.replyPong(m, src.String())
	case proto.MsgPong:
		n.touch(m.From.ID, m.From.Addr)
		n.rpcs.Resolve(m.RPC, m)
	default:
		// later: FIND_NODE, etc.
	}
}

func (n *Node) touch(hexID, addr string) {
	ni, err := id.FromHex(hexID)
	if err != nil { return }
	n.mu.Lock()
	defer n.mu.Unlock()
	// move-to-front or insert
	idx := -1
	for i, c := range n.contacts {
		if c.addr == addr || c.id == ni {
			idx = i
			break
		}
	}
	if idx >= 0 {
		c := n.contacts[idx]
		c.lastSeen = time.Now()
		c.addr = addr
		n.contacts = append([]contact{c}, append(n.contacts[:idx], n.contacts[idx+1:]...)...)
	} else {
		n.contacts = append([]contact{{id: ni, addr: addr, lastSeen: time.Now()}}, n.contacts...)
		if len(n.contacts) > 2000 { // arbitrary cap
			n.contacts = n.contacts[:2000]
		}
	}
}

func (n *Node) replyPong(req proto.Message, to string) {
	resp := proto.Message{
		Type: proto.MsgPong,
		RPC:  req.RPC,
		From: proto.Contact{ID: id.ToHex(n.selfID), Addr: n.tr.Addr()},
		To:   to,
	}
	_ = n.tr.Send(to, resp)
}

func (n *Node) Ping(to string, peerID id.NodeID) (bool, error) {
	rpc := id.NewRandomRPC()
	ch := n.rpcs.Add(id.ToHex(rpc))
	msg := proto.Message{
		Type: proto.MsgPing,
		RPC:  id.ToHex(rpc),
		From: proto.Contact{ID: id.ToHex(n.selfID), Addr: n.tr.Addr()},
		To:   to,
	}
	if err := n.tr.Send(to, msg); err != nil {
		return false, err
	}
	select {
	case <-ch:
		return true, nil
	case <-time.After(800 * time.Millisecond):
		return false, nil
	}
}