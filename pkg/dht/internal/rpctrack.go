// pkg/dht/internal/rpctrack.go
package internal

import (
	"sync"
	"github.com/BrandonChongWenJun/D7024e-tutorial/pkg/proto"
)

type Waiter struct {
	ch chan proto.Message
}

type RPCTrack struct {
	mu sync.Mutex
	m  map[string]*Waiter // rpc hex -> waiter
}

func NewRPCTrack() *RPCTrack { return &RPCTrack{m: make(map[string]*Waiter)} }

func (t *RPCTrack) Add(rpc string) chan proto.Message {
	t.mu.Lock()
	defer t.mu.Unlock()
	w := &Waiter{ch: make(chan proto.Message, 1)}
	t.m[rpc] = w
	return w.ch
}

func (t *RPCTrack) Resolve(rpc string, msg proto.Message) {
	t.mu.Lock()
	w := t.m[rpc]
	delete(t.m, rpc)
	t.mu.Unlock()
	if w != nil {
		w.ch <- msg
		close(w.ch)
	}
}