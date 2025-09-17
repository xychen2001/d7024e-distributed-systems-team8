package network

import (
	"encoding/json"
	"fmt"
	"github.com/xychen2001/d7024e-distributed-systems-team8/pkg/dht"
)

// MessageType represents the type of a Kademlia message.
type MessageType int

const (
	PING MessageType = iota
	PONG
	FIND_NODE
	STORE
	FIND_VALUE
)

// Message represents a Kademlia message.
type Message struct {
	RPCID    *dht.KademliaID
	SenderID *dht.KademliaID
	Type     MessageType
	Payload  []byte
}

// Serialize converts a Message to a byte slice for network transmission.
func (m *Message) Serialize() ([]byte, error) {
	return json.Marshal(m)
}

// Deserialize converts a byte slice back to a Message.
func Deserialize(data []byte) (*Message, error) {
	var msg Message
	err := json.Unmarshal(data, &msg)
	return &msg, err
}

// String returns a string representation of the MessageType.
func (mt MessageType) String() string {
	switch mt {
	case PING:
		return "PING"
	case PONG:
		return "PONG"
	case FIND_NODE:
		return "FIND_NODE"
	case STORE:
		return "STORE"
	case FIND_VALUE:
		return "FIND_VALUE"
	default:
		return fmt.Sprintf("UNKNOWN(%d)", mt)
	}
}
