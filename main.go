package main

import (
	"fmt"
	"os"

	"github.com/pusher/cli/api"
	"github.com/pusher/cli/commands/auth"
	"github.com/pusher/cli/commands/channels"
	"github.com/pusher/cli/commands/cli"
	"github.com/pusher/cli/config"
	"github.com/spf13/cobra"
)

func main() {
	config.Init()
	var rootCmd = &cobra.Command{Use: "pusher", Short: "A CLI for your Pusher account. Find out more at https://pusher.com"}

	var Apps = &cobra.Command{Use: "apps",
		Short: "Manage your Channels Apps"}
	funcCmd, err := channels.NewFunctionsCommand(api.NewPusherApi())
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not execute command: %s\n", err.Error())
		os.Exit(1)
		return
	}
	Apps.AddCommand(channels.Apps, channels.Tokens, channels.Subscribe, channels.Trigger, channels.ListChannels, channels.ChannelInfo, funcCmd)

	var Generate = &cobra.Command{Use: "generate",
		Short: "Generate a Channels client, server, or Authorisation server"}
	Generate.AddCommand(channels.GenerateClient, channels.GenerateServer, channels.LocalAuthServer)

	var Channels = &cobra.Command{Use: "channels",
		Short: "Commands related to the Channels product"}
	Channels.AddCommand(Generate, Apps)

	rootCmd.AddCommand(Channels)
	rootCmd.AddCommand(auth.Login, auth.Logout)
	rootCmd.AddCommand(cli.Version)
	err = rootCmd.Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not execute command: %s\n", err.Error())
		os.Exit(1)
		return
	}

}
