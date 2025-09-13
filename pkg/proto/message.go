// pkg/proto/message.go
package proto

//import "github.com/BrandonChongWenJun/D7024e-tutorial/pkg/id"

type MsgType string

const (
	MsgPing       MsgType = "PING"
	MsgPong       MsgType = "PONG"
	MsgFindNode   MsgType = "FIND_NODE"
	MsgNodes      MsgType = "NODES"
	MsgStore      MsgType = "STORE"
	MsgFindValue  MsgType = "FIND_VALUE"
	MsgValue      MsgType = "VALUE"
)

type Contact struct {
	ID   string // hex NodeID
	Addr string // "ip:port"
}

type Message struct {
	Type  MsgType   `json:"type"`
	RPC   string    `json:"rpc"`   // hex id.RPCID
	From  Contact   `json:"from"`
	To    string    `json:"to"`    // addr string
	// Payloads:
	Target string   `json:"target,omitempty"` // hex NodeID (FIND_NODE / FIND_VALUE)
	Nodes  []Contact `json:"nodes,omitempty"` // NODES
	Key    string   `json:"key,omitempty"`    // hex NodeID
	Value  string   `json:"value,omitempty"`  // UTF-8 value (<=255B)
}