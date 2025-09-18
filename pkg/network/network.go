// pkg/network/network.go
package network

import (
	"encoding/json"
	"errors"
	"log"
	"net"
	"sync"
	"time"

	"github.com/xychen2001/d7024e-distributed-systems-team8/pkg/dht"
)

const rpcTimeout = 5 * time.Second

// Network handles the UDP communication between nodes.
type Network struct {
	NodeID           *dht.KademliaID
	ListenAddr       string
	conn             *net.UDPConn
	routingTable     *dht.RoutingTable
	mutex            sync.RWMutex
	pendingResponses map[dht.KademliaID]chan *Message
}

// NewNetwork creates a new Network instance.
func NewNetwork(nodeID *dht.KademliaID, rt *dht.RoutingTable, listenAddr string) *Network {
	return &Network{
		NodeID:           nodeID,
		ListenAddr:       listenAddr,
		routingTable:     rt,
		pendingResponses: make(map[dht.KademliaID]chan *Message),
	}
}

// Listen starts the UDP listener for incoming messages.
func (n *Network) Listen() {
	addr, err := net.ResolveUDPAddr("udp", n.ListenAddr)
	if err != nil {
		log.Fatalf("Failed to resolve UDP address: %v", err)
	}

	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		log.Fatalf("Failed to listen on UDP address: %v", err)
	}
	n.conn = conn
	log.Printf("Listening on %s\n", n.ListenAddr)

	go func() {
		defer conn.Close()
		buffer := make([]byte, 4096) // Increased buffer size for larger payloads
		for {
			length, remote, err := conn.ReadFromUDP(buffer)
			if err != nil {
				log.Printf("Error reading from UDP: %v", err)
				continue
			}
			go n.handleMessage(buffer[:length], remote)
		}
	}()
}

// handleMessage deserializes and processes an incoming message.
func (n *Network) handleMessage(data []byte, remote *net.UDPAddr) {
	msg, err := Deserialize(data)
	if err != nil {
		log.Printf("Error deserializing message from %s: %v", remote, err)
		return
	}

	// Check if this is a response to a pending RPC
	n.mutex.RLock()
	responseChan, isResponse := n.pendingResponses[*msg.RPCID]
	n.mutex.RUnlock()

	if isResponse {
		responseChan <- msg
		return
	}

	// Otherwise, handle it as a new request
	log.Printf("Received %s from %s", msg.Type, remote)

	switch msg.Type {
	case PING:
		pongMsg := Message{
			RPCID:    msg.RPCID,
			SenderID: n.NodeID,
			Type:     PONG,
		}
		n.sendMessage(&pongMsg, remote)
	case FIND_NODE:
		var targetID dht.KademliaID
		if err := json.Unmarshal(msg.Payload, &targetID); err != nil {
			log.Printf("Failed to unmarshal FIND_NODE payload: %v", err)
			return
		}
		closestContacts := n.routingTable.FindClosestContacts(&targetID, dht.BucketSize)
		payload, err := json.Marshal(closestContacts)
		if err != nil {
			log.Printf("Failed to marshal closest contacts: %v", err)
			return
		}
		responseMsg := Message{
			RPCID:    msg.RPCID,
			SenderID: n.NodeID,
			Type:     FIND_NODE, // Response type is the same
			Payload:  payload,
		}
		n.sendMessage(&responseMsg, remote)
	default:
		log.Printf("Received unknown message type %d from %s", msg.Type, remote)
	}
}

// sendMessage serializes and sends a message to a remote address.
func (n *Network) sendMessage(msg *Message, remote *net.UDPAddr) {
	data, err := msg.Serialize()
	if err != nil {
		log.Printf("Error serializing message for %s: %v", remote, err)
		return
	}

	_, err = n.conn.WriteToUDP(data, remote)
	if err != nil {
		log.Printf("Error sending message to %s: %v", remote, err)
	}
}

// FindNode sends a FIND_NODE request and waits for a response.
func (n *Network) FindNode(contact *dht.Contact, target *dht.KademliaID) ([]dht.Contact, error) {
	rpcID := dht.NewRandomKademliaID()
	payload, err := json.Marshal(target)
	if err != nil {
		return nil, err
	}

	requestMsg := &Message{
		RPCID:    rpcID,
		SenderID: n.NodeID,
		Type:     FIND_NODE,
		Payload:  payload,
	}

	responseChan := make(chan *Message, 1)
	n.mutex.Lock()
	n.pendingResponses[*rpcID] = responseChan
	n.mutex.Unlock()

	defer func() {
		n.mutex.Lock()
		delete(n.pendingResponses, *rpcID)
		n.mutex.Unlock()
	}()

	remoteAddr, err := net.ResolveUDPAddr("udp", contact.Address)
	if err != nil {
		return nil, err
	}

	n.sendMessage(requestMsg, remoteAddr)

	select {
	case responseMsg := <-responseChan:
		var contacts []dht.Contact
		if err := json.Unmarshal(responseMsg.Payload, &contacts); err != nil {
			return nil, err
		}
		return contacts, nil
	case <-time.After(rpcTimeout):
		return nil, errors.New("rpc timeout")
	}
}

// Ping sends a PING request and waits for a PONG response.
func (n *Network) Ping(contact *dht.Contact) error {
	rpcID := dht.NewRandomKademliaID()

	requestMsg := &Message{
		RPCID:    rpcID,
		SenderID: n.NodeID,
		Type:     PING,
	}

	responseChan := make(chan *Message, 1)
	n.mutex.Lock()
	n.pendingResponses[*rpcID] = responseChan
	n.mutex.Unlock()

	defer func() {
		n.mutex.Lock()
		delete(n.pendingResponses, *rpcID)
		n.mutex.Unlock()
	}()

	remoteAddr, err := net.ResolveUDPAddr("udp", contact.Address)
	if err != nil {
		return err
	}

	n.sendMessage(requestMsg, remoteAddr)

	select {
	case responseMsg := <-responseChan:
		if responseMsg.Type == PONG {
			// The contact ID should be updated from the PONG response
			contact.ID = responseMsg.SenderID
			return nil
		} else {
			return errors.New("invalid response type for ping")
		}
	case <-time.After(rpcTimeout):
		return errors.New("rpc timeout for ping")
	}
}