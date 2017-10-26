package main

import (
	"github.com/pusher/cli/commands/auth"
	"github.com/pusher/cli/commands/channels"
	"github.com/spf13/cobra"
)

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
	rootCmd.Execute()
}
