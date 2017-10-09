package main

import (
	"github.com/pusher/pusher-cli/commands/auth"
	"github.com/pusher/pusher-cli/commands/ddn"
	"github.com/spf13/cobra"
)

func main() {
	var rootCmd = &cobra.Command{Use: "pusher"}

	var Auth = &cobra.Command{Use: "auth",
		Short: "Add your API key to access the Pusher Platform"}

	Auth.AddCommand(auth.Login, auth.Logout, auth.Status)
	var Apps = &cobra.Command{Use: "apps",
		Short: "Manage your Pusher Apps"}
	Apps.AddCommand(ddn.Apps, ddn.Tokens, ddn.Subscribe, ddn.Trigger)

	var Generate = &cobra.Command{Use: "generate",
		Short: "Generate a Pusher client, server, or Authorisation server"}
	Generate.AddCommand(ddn.GenerateClient, ddn.GenerateServer, ddn.LocalAuthServer)

	var DDN = &cobra.Command{Use: "ddn",
		Short: "Commands related to our PubSub product."}
	DDN.AddCommand(Generate, Apps)

	rootCmd.AddCommand(Auth)
	rootCmd.AddCommand(DDN)
	rootCmd.Execute()
}
