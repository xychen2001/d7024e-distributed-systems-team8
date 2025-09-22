package cli

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"
	"github.com/xychen2001/d7024e-distributed-systems-team8/pkg/dht"
	"github.com/xychen2001/d7024e-distributed-systems-team8/pkg/network"
)

var bootstrapAddress string
var port int

func init() {
	startCmd.Flags().StringVarP(&bootstrapAddress, "bootstrap", "b", "", "Address of a bootstrap node to join the network")
	startCmd.Flags().IntVarP(&port, "port", "p", 8080, "Port to listen on")
	rootCmd.AddCommand(startCmd)
}

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Starts a Kademlia node",
	Long:  `Starts a Kademlia node, which will begin listening for incoming UDP messages.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Construct the listen address from the port flag.
		listenAddr := fmt.Sprintf("127.0.0.1:%d", port)

		nodeID := dht.NewRandomKademliaID()
		log.Printf("Starting node with ID %s on %s", nodeID, listenAddr)

		// Create the node's own contact details.
		me := dht.NewContact(nodeID, listenAddr)

		// Create a new routing table.
		rt := dht.NewRoutingTable(me)

		// Create the network layer.
		net := network.NewNetwork(nodeID, rt, listenAddr)

		// Start the network listener.
		net.Listen()

		// Create the Kademlia instance.
		kademlia := dht.NewKademlia(rt, net)

		// If a bootstrap address is provided, join the network.
		if bootstrapAddress != "" {
			go func() {
				log.Printf("Joining network via bootstrap node at %s...", bootstrapAddress)

				// Create a temporary contact for the bootstrap node (ID is unknown)
				bootstrapContact := dht.NewContact(dht.NewRandomKademliaID(), bootstrapAddress)

				// Ping the bootstrap node to get its real ID
				if err := net.Ping(&bootstrapContact); err != nil {
					log.Printf("Failed to ping bootstrap node: %v", err)
					return
				}
				// The bootstrapContact ID is now updated from the PONG response.
				log.Printf("Successfully contacted bootstrap node with ID %s", bootstrapContact.ID)

				// Add the now-known bootstrap contact to the routing table
				rt.AddContact(bootstrapContact, net)

				// Perform a lookup for our own ID to populate the routing table.
				kademlia.LookupContact(me.ID)
				log.Println("Bootstrap process finished. Node is now part of the network.")
			}()
		}

		// Block forever to keep the listener running
		select {}
	},
}
