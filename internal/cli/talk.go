package cli

import (
	"github.com/xychen2001/d7024e-distributed-systems-team8/pkg/helloworld"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(TalkCmd)
}

var TalkCmd = &cobra.Command{
	Use:   "talk",
	Short: "Say something",
	Long:  "Say something",
	Run: func(cmd *cobra.Command, args []string) {
		hellworld := helloworld.NewHelloWorld()
		hellworld.Talk()
	},
}
