// pkg/dht/node.go
package dht

import (
	"context"
	"sync"
	"time"

	"github.com/BrandonChongWenJun/D7024e-tutorial/pkg/id"
	//"github.com/BrandonChongWenJun/D7024e-tutorial/pkg/proto"
	"github.com/BrandonChongWenJun/D7024e-tutorial/pkg/transport"
	"github.com/BrandonChongWenJun/D7024e-tutorial/pkg/dht/internal"
)

type contact struct {
	id       id.NodeID
	addr     string
	lastSeen time.Time
}

type Node struct {
	cfg     Config
	selfID  id.NodeID
	tr      transport.Transport

	mu      sync.RWMutex
	rt      *RoutingTable

	rpcs *internal.RPCTrack
}

func New(ctx context.Context, cfg Config, tr transport.Transport) *Node {
	n := &Node{
		cfg: cfg,
		tr:  tr,
		rpcs: internal.NewRPCTrack(),
	}
	// Use SHA1 of addr as deterministic NodeID for now
	n.selfID = id.HashSHA1([]byte(cfg.SelfAddr))
	// Initialize routing table (K-buckets)
	if n.cfg.K == 0 {
		n.cfg.K = 20
	}
	n.rt = newRoutingTable(n.selfID, n.cfg.K)
	tr.SetHandler(n.handle)
	go tr.Start(ctx)
	return n
}

func (n *Node) Self() (id.NodeID, string) { return n.selfID, n.tr.Addr() }