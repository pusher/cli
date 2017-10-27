package main

import (
	"github.com/pusher/cli/commands/auth"
	"github.com/pusher/cli/commands/channels"
	"github.com/pusher/cli/commands/cli"
	"github.com/spf13/cobra"
)

const pusherLogo = `
            ▄▓▓▓▓▓▓▓▄.        
      .▄█,▓▓▓▓▓▓▀  ▄#▓█▓▓╥    
    ╓▓▓ ▄▓▓▓▓▓▀ ,▓▓▓▓▓▓▓▓▓▓▓▄ 
   ▐▓▓ ;▓▓▓▓▓▌ ▓▓▓▓▓▓▀▄▄ ▀▓▓▓▓
   ▓▓▌ ▓▓▓▓▓▓ ▓▓▓▓▓▓  ▓▌  ▀▓▓▓
  ╟▓▓▌ ▓▓▓▓▓▌ ▓▓▓▓▓  ▓▓   ▐▓▓▌
  ▓▓▓▌ ▓▓▓▓▓▌ ▓▓▓▓▓ ▀   .▓▓▀ 
  ▓▓▓▌ ▓▓▓▓▓▌ ▓▓▓▓▌    .▓▓▓^  
  ▐▓▓▌ ▓▓▓▓▓▓ ▓▓▓▓▌ ,▓▓▓▀     
   ▓▓▓ ▓▓▓▓▓▓ ▓▓▓▓▓ ▀        
   ▓▓▓ ╟▓▓▓▓▓▌╙▓▓▓▓           
   ▐▓▓▓ ▓▓▓▓▓▓ ▓▓▓▓▓  ,,▄     
    ▓▓▓ ▀▓▓▓▓▓▄▀▓▓▓▓▓▓▓▓      
     ▓▓▓ ▓▓▓▓▓▓ ▀▀▀▀ ▓▓       
ⁿ█▓▓▓▓▓█═▐▓▓▓▓▓▓▓▓▓▀▐▓        
   ▀▓▄ █▓▓▓▓▓▓▓▓▓▓▀┌▓        
     ▀▓⌐▀▓▓▓▓▓▓▓▀ █▀         
            ▀▓▓▓▀             
               ▀        
`

func main() {
	var rootCmd = &cobra.Command{Use: "pusher", Short: pusherLogo}
	var Apps = &cobra.Command{Use: "apps",
		Short: "Manage your Pusher Apps"}
	Apps.AddCommand(channels.Apps, channels.Tokens, channels.Subscribe, channels.Trigger)

	var Generate = &cobra.Command{Use: "generate",
		Short: "Generate a Pusher client, server, or Authorisation server"}

	var Channels = &cobra.Command{Use: "channels",
		Short: "Commands related to Pusher Channels"}

	Generate.AddCommand(channels.GenerateClient, channels.GenerateServer, channels.LocalAuthServer)
	Channels.AddCommand(Generate, Apps)
	rootCmd.AddCommand(Channels)
	rootCmd.AddCommand(auth.Login, auth.Logout)
	rootCmd.AddCommand(cli.Version)
	rootCmd.Execute()
}
