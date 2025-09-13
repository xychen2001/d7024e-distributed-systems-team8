// pkg/dht/pingpong_test.go
package dht_test

import (
	"context"
	"testing"

	"github.com/BrandonChongWenJun/D7024e-tutorial/pkg/dht"
	"github.com/BrandonChongWenJun/D7024e-tutorial/pkg/transport/udp"
)

func TestPingPong(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	a := udp.MustNew("127.0.0.1:0")
	b := udp.MustNew("127.0.0.1:0")

	na := dht.New(ctx, dht.Config{SelfAddr: a.Addr()}, a)
	nb := dht.New(ctx, dht.Config{SelfAddr: b.Addr()}, b)
	peerID, _ := nb.Self()
	ok, err := na.Ping(b.Addr(), peerID) // or derive peerID from addr in your API
	if err != nil || !ok {
		t.Fatalf("ping failed: %v", err)
	}
}