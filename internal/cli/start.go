package cli

import (
	"log"

	"github.com/spf13/cobra"
	"github.com/xychen2001/d7024e-distributed-systems-team8/pkg/dht"
	"github.com/xychen2001/d7024e-distributed-systems-team8/pkg/network"
)

func init() {
	rootCmd.AddCommand(startCmd)
}

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Starts a Kademlia node",
	Long:  `Starts a Kademlia node, which will begin listening for incoming UDP messages.`,
	Run: func(cmd *cobra.Command, args []string) {
		// TODO: Get listen address from flags
		listenAddr := "0.0.0.0:8080"

		nodeID := dht.NewRandomKademliaID()
		log.Printf("Starting node with ID %s on %s", nodeID, listenAddr)

		net := network.NewNetwork(nodeID, listenAddr)
		net.Listen()

		// Block forever to keep the listener running
		select {}
	},
}
