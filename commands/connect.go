package commands

import (
	"github.com/spf13/cobra"
)

var Connect = &cobra.Command{
	Use:   "connect",
	Short: "Open WebSocket connection to Pusher, allowing subscription to channels",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
	},
}
