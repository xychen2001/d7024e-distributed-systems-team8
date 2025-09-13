// pkg/transport/transport.go
package transport

import (
	"context"
	"net"

	"github.com/BrandonChongWenJun/D7024e-tutorial/pkg/proto"
)

type Handler func(msg proto.Message, src *net.UDPAddr)

type Transport interface {
	Start(ctx context.Context) error
	Send(to string, msg proto.Message) error
	Addr() string
	SetHandler(h Handler)
	Close() error
}