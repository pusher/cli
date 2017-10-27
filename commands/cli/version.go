package cli

import (
	"fmt"

	"github.com/pusher/cli/config"
	"github.com/spf13/cobra"
)

// Version lists the version of Pusher CLI
var Version = &cobra.Command{
	Use:   "version",
	Short: "Lists the version of Pusher CLI",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Pusher CLI: " + config.GetVersion())
	},
}
