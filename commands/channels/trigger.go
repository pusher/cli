package channels

import (
	"fmt"
	"os"

	"github.com/pusher/cli/api"
	"github.com/pusher/cli/commands"
	"github.com/pusher/pusher-http-go"
	"github.com/spf13/cobra"
)

// Trigger allows the user to trigger an event on a particular channel.
var Trigger = &cobra.Command{
	Use:   "trigger",
	Short: "Trigger an event on a Channels app",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {

		if commands.AppID == "" {
			fmt.Fprintf(os.Stderr, "Please supply --app-id\n")
			os.Exit(1)
			return
		}

		if commands.ChannelName == "" {
			fmt.Fprintf(os.Stderr, "Please supply --channel\n")
			os.Exit(1)
			return
		}

		if commands.EventName == "" {
			fmt.Fprintf(os.Stderr, "Please supply --event\n")
			os.Exit(1)
			return
		}

		if commands.Message == "" {
			fmt.Fprintf(os.Stderr, "Please supply --message\n")
			os.Exit(1)
			return
		}

		app, err := api.GetApp(commands.AppID)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Could not get app the app: %s\n", err.Error())
			os.Exit(1)
			return
		}

		token, err := api.GetToken(commands.AppID)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Could not get app token: %s\n", err.Error())
			os.Exit(1)
			return
		}

		client := pusher.Client{
			AppId:   commands.AppID,
			Key:     token.Key,
			Secret:  token.Secret,
			Cluster: app.Cluster + ".staging", // app.Cluster,
		}

		_, err = client.Trigger(commands.ChannelName, commands.EventName, commands.Message)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Could not trigger: %s\n", err.Error())
			return
		}
	},
}

func init() {
	Trigger.PersistentFlags().StringVar(&commands.AppID, "app-id", "", "Channels App ID")
	Trigger.PersistentFlags().StringVar(&commands.ChannelName, "channel", "", "Channel name")
	Trigger.PersistentFlags().StringVar(&commands.EventName, "event", "", "Event name")
	Trigger.PersistentFlags().StringVar(&commands.Message, "message", "", "Message")
}
