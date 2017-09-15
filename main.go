package main

import (
	"github.com/pusher/pusher-cli/commands"
	"github.com/spf13/cobra"
)

func main() {
	var rootCmd = &cobra.Command{Use: "pusher"}
	rootCmd.AddCommand(commands.Login)
	rootCmd.AddCommand(commands.Logout)
	rootCmd.AddCommand(commands.Apps)
	rootCmd.AddCommand(commands.Trigger)
	rootCmd.AddCommand(commands.Subscribe)
	rootCmd.AddCommand(commands.Generate)
	rootCmd.Execute()
}
