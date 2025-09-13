// pkg/transport/udp/udp.go
package udp

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"sync"
	"time"

	"github.com/BrandonChongWenJun/D7024e-tutorial/pkg/proto"
	"github.com/BrandonChongWenJun/D7024e-tutorial/pkg/transport"
)

type UDP struct {
	conn    *net.UDPConn
	addr    string
	handler transport.Handler
	mu      sync.RWMutex
}

func New(bind string) (*UDP, error) {
	udpAddr, err := net.ResolveUDPAddr("udp", bind)
	if err != nil {
		return nil, err
	}
	c, err := net.ListenUDP("udp", udpAddr)
	if err != nil {
		return nil, err
	}
	return &UDP{conn: c, addr: c.LocalAddr().String()}, nil
}

func (u *UDP) Addr() string { return u.addr }

func (u *UDP) SetHandler(h transport.Handler) {
	u.mu.Lock()
	u.handler = h
	u.mu.Unlock()
}

func (u *UDP) Start(ctx context.Context) error {
	buf := make([]byte, 64*1024)
	for {
		u.conn.SetReadDeadline(time.Now().Add(2 * time.Second))
		n, src, err := u.conn.ReadFromUDP(buf)
		if ne, ok := err.(net.Error); ok && ne.Timeout() {
			select {
			case <-ctx.Done():
				return nil
			default:
				continue
			}
		}
		if err != nil {
			select {
			case <-ctx.Done():
				return nil
			default:
				continue
			}
		}
		var m proto.Message
		if err := json.Unmarshal(buf[:n], &m); err != nil {
			continue
		}
		u.mu.RLock()
		h := u.handler
		u.mu.RUnlock()
		if h != nil {
			go h(m, src)
		}
	}
}

func (u *UDP) Send(to string, msg proto.Message) error {
	b, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	addr, err := net.ResolveUDPAddr("udp", to)
	if err != nil {
		return err
	}
	_, err = u.conn.WriteToUDP(b, addr)
	return err
}

func (u *UDP) Close() error { return u.conn.Close() }

func MustNew(bind string) *UDP {
	udp, err := New(bind)
	if err != nil {
		panic(fmt.Errorf("udp: %w", err))
	}
	return udp
}