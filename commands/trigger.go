package commands

import (
	"fmt"
	"os"

	"github.com/pusher/pusher-cli/api"
	"github.com/pusher/pusher-http-go"
	"github.com/spf13/cobra"
)

var Trigger = &cobra.Command{
	Use:   "trigger",
	Short: "Trigger an event on a Pusher app",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {

		if appId == "" {
			fmt.Fprintf(os.Stderr, "Please supply --app-id\n")
			os.Exit(1)
			return
		}

		if channelName == "" {
			fmt.Fprintf(os.Stderr, "Please supply --channel\n")
			os.Exit(1)
			return
		}

		if eventName == "" {
			fmt.Fprintf(os.Stderr, "Please supply --event\n")
			os.Exit(1)
			return
		}

		if message == "" {
			fmt.Fprintf(os.Stderr, "Please supply --message\n")
			os.Exit(1)
			return
		}

		app, err := api.GetApp(appId)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Could not get the app: %s\n", err.Error())
			os.Exit(1)
			return
		}

		token, err := api.GetToken(appId)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Could not get token: %s\n", err.Error())
			os.Exit(1)
			return
		}

		client := pusher.Client{
			AppId:   appId,
			Key:     token.Key,
			Secret:  token.Secret,
			Cluster: app.Cluster + ".staging", // app.Cluster,
		}

		_, err = client.Trigger(channelName, eventName, message)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Could not trigger: %s\n", err.Error())
			return
		}
	},
}

func init() {
	Trigger.PersistentFlags().StringVar(&appId, "app-id", "", "Pusher App ID")
	Trigger.PersistentFlags().StringVar(&channelName, "channel", "", "Channel name")
	Trigger.PersistentFlags().StringVar(&eventName, "event", "", "Event name")
	Trigger.PersistentFlags().StringVar(&message, "message", "", "Message")
}
