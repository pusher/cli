package commands

import (
	"fmt"
	"github.com/spf13/cobra"
)

var Generate = &cobra.Command{
	Use:   "generate",
	Short: "Generate a starter project which uses a Pusher app",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("TODO Generating something")
	},
}
