package network

import (
	"log"
	"net"

	"github.com/xychen2001/d7024e-distributed-systems-team8/pkg/dht"
)

// Network handles the UDP communication between nodes.
type Network struct {
	NodeID     *dht.KademliaID
	ListenAddr string
	conn       *net.UDPConn
}

// NewNetwork creates a new Network instance.
func NewNetwork(nodeID *dht.KademliaID, listenAddr string) *Network {
	return &Network{
		NodeID:     nodeID,
		ListenAddr: listenAddr,
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
		buffer := make([]byte, 2048) // Buffer for incoming messages
		for {
			length, remote, err := conn.ReadFromUDP(buffer)
			if err != nil {
				log.Printf("Error reading from UDP: %v", err)
				continue
			}
			// Handle each incoming message in a new goroutine for concurrency
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

	log.Printf("Received %s from %s", msg.Type, remote)

	switch msg.Type {
	case PING:
		// Respond with a PONG message
		pongMsg := Message{
			RPCID:    msg.RPCID, // Echo the RPCID back
			SenderID: n.NodeID,
			Type:     PONG,
		}
		n.sendMessage(&pongMsg, remote)
	case PONG:
		// In a real scenario, we would handle the PONG, e.g., by notifying a waiting process.
		log.Printf("Received PONG from %s", remote)
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
