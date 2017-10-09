package commands

import (
	"fmt"

	"github.com/pusher/pusher-cli/commands/auth"
	"github.com/spf13/cobra"
)

// CheckAuth checks the stored API key for validity
var CheckAuth = &cobra.Command{
	Use:   "check-auth",
	Short: "Checks the stored API key for validity.",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {

		if auth.APIKeyValid() {
			fmt.Println("API Key is valid.")
		} else {
			fmt.Println("API Key is NOT valid.")
		}
	},
}
