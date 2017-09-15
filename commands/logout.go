package commands

import (
	"fmt"
	"os"

	"github.com/pusher/pusher-cli/config"
	"github.com/spf13/cobra"
)

var Logout = &cobra.Command{
	Use:   "logout",
	Short: "Remove Pusher account credentials from this computer",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {

		deleteErr := config.Delete()
		if deleteErr == nil {
			fmt.Println("Removed Pusher account credentials.")
		} else {
			fmt.Println("Failed to remove the local configuration.")
			os.Exit(1)
		}
	},
}
