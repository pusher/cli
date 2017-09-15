package commands

import (
	"fmt"
	"github.com/spf13/cobra"
)

var Publish = &cobra.Command{
	Use:   "publish [app-id] [channel] [event] [message]",
	Short: "Trigger an event on a Pusher app",
	Args: cobra.ExactArgs(4),
	Run: func(cmd *cobra.Command, args []string) {
		appId := args[0]
		channelName := args[1]
		eventName := args[2]
		message := args[3]
		fmt.Printf("Triggered event: %s %s %s %s\n", appId, channelName, eventName, message)
	},
}
