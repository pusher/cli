package main

import (
	"github.com/spf13/cobra"
	"github.com/pusher/pusher-cli/commands"
)

func main() {
	var rootCmd = &cobra.Command{Use: "pusher"}
	rootCmd.AddCommand(commands.Login)
	rootCmd.AddCommand(commands.Logout)
	rootCmd.AddCommand(commands.Publish)
	rootCmd.Execute()
}
