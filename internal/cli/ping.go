package cli

import (
	"fmt"
	"log"
	"net"
	"time"

	"github.com/spf13/cobra"
	"github.com/xychen2001/d7024e-distributed-systems-team8/pkg/dht"
	"github.com/xychen2001/d7024e-distributed-systems-team8/pkg/network"
)

func init() {
	rootCmd.AddCommand(pingCmd)
}

var pingCmd = &cobra.Command{
	Use:   "ping [address]",
	Short: "Sends a PING message to a node to test connectivity.",
	Long:  `Sends a UDP PING message to a specified node address and waits for a PONG response. This is used to verify that a node is alive and reachable.`, 
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		targetAddr := args[0]
		log.Printf("Pinging %s...", targetAddr)

		// Use net.Dial to get a connection that can write and read.
		// This will also resolve the address for us.
		conn, err := net.Dial("udp", targetAddr)
		if err != nil {
			log.Fatalf("Failed to connect to target: %v", err)
		}
		defer conn.Close()

		// Create a PING message with a unique RPC ID.
		rpcID := dht.NewRandomKademliaID()
		senderID := dht.NewRandomKademliaID() // This is a dummy ID for the ping command.
		pingMsg := network.Message{
			RPCID:    rpcID,
			SenderID: senderID,
			Type:     network.PING,
		}

		data, err := pingMsg.Serialize()
		if err != nil {
			log.Fatalf("Failed to serialize PING message: %v", err)
		}

		_, err = conn.Write(data)
		if err != nil {
			log.Fatalf("Failed to send PING message: %v", err)
		}

		log.Println("PING sent. Waiting for PONG...")

		// Set a deadline for the read operation to avoid waiting forever.
		conn.SetReadDeadline(time.Now().Add(5 * time.Second))

		buffer := make([]byte, 2048)
		n, err := conn.Read(buffer)
		if err != nil {
			log.Fatalf("Failed to receive PONG: %v", err)
		}

		pongMsg, err := network.Deserialize(buffer[:n])
		if err != nil {
			log.Fatalf("Failed to deserialize PONG message: %v", err)
		}

		// Verify the response is a PONG and the RPC ID matches.
		if pongMsg.Type == network.PONG && pongMsg.RPCID.Equals(rpcID) {
			fmt.Printf("\nSuccess! PONG received from node %s\n", pongMsg.SenderID)
		} else {
			fmt.Printf("\nError: Received an unexpected message of type %s\n", pongMsg.Type)
		}
	},
}
