package main

import (
	"fmt"

	"github.com/pusher/cli/commands/auth"
	"github.com/pusher/cli/commands/channels"
	"github.com/spf13/cobra"
)

var VERSION = "master"

// Version lists the version of Pusher CLI
var Version = &cobra.Command{
	Use:   "version",
	Short: "Lists the version of Pusher CLI",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Pusher CLI: " + VERSION)
	},
}

func main() {
	var rootCmd = &cobra.Command{Use: "pusher"}
	var Apps = &cobra.Command{Use: "apps",
		Short: "Manage your Pusher Apps"}
	Apps.AddCommand(channels.Apps, channels.Tokens, channels.Subscribe, channels.Trigger)

	var Generate = &cobra.Command{Use: "generate",
		Short: "Generate a Pusher client, server, or Authorisation server"}
	Generate.AddCommand(channels.GenerateClient, channels.GenerateServer, channels.LocalAuthServer)
	rootCmd.AddCommand(Generate, Apps)
	rootCmd.AddCommand(auth.Login, auth.Logout)
	rootCmd.AddCommand(Version)
	rootCmd.Execute()
}
