package main

import (
	"github.com/pusher/pusher-cli/commands/auth"
	"github.com/pusher/pusher-cli/commands/ddn"
	"github.com/spf13/cobra"
)

func main() {
	var rootCmd = &cobra.Command{Use: "pusher"}
	var Apps = &cobra.Command{Use: "apps",
		Short: "Manage your Pusher Apps"}
	Apps.AddCommand(ddn.Apps, ddn.Tokens, ddn.Subscribe, ddn.Trigger)

	var Generate = &cobra.Command{Use: "generate",
		Short: "Generate a Pusher client, server, or Authorisation server"}
	Generate.AddCommand(ddn.GenerateClient, ddn.GenerateServer, ddn.LocalAuthServer)
	rootCmd.AddCommand(Generate, Apps)
	rootCmd.AddCommand(auth.Login, auth.Logout)
	rootCmd.Execute()
}
