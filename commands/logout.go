package commands

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

var Logout = &cobra.Command{
	Use:   "logout",
	Short: "Remove Pusher account credentials from this computer",
	Args: cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Print: " + strings.Join(args, " "))
	},
}
