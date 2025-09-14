package cli

import (
	"fmt"

	"github.com/xychen2001/d7024e-distributed-systems-team8/pkg/build"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(VersionCmd)
}

var VersionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version",
	Long:  "Print the version of the application",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(build.BuildVersion)
		fmt.Println(build.BuildTime)
	},
}
