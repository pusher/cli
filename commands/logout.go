package commands

import (
	"fmt"
	"github.com/spf13/cobra"
)

var Logout = &cobra.Command{
	Use:   "logout",
	Short: "Remove Pusher account credentials from this computer",
	Args: cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Removed Pusher account credentials. Run `pusher login` to login again.")
	},
}
