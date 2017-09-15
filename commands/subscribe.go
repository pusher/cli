package commands

import (
	"fmt"
	"os"
	"time"

	"github.com/fatih/color"
	"github.com/pusher-community/pusher-websocket-go"
	"github.com/pusher/pusher-cli/api"
	"github.com/spf13/cobra"
)

var Subscribe = &cobra.Command{
	Use:   "subscribe",
	Short: "Subscribe to a Pusher channel",
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

		pusher.New(token.Key)
		client := pusher.NewWithConfig(pusher.ClientConfig{
			Scheme: "wss",
			Host:   "ws-" + app.Cluster + ".staging.pusher.com",
			Port:   "443",
			Key:    token.Key,
			Secret: token.Secret,
		})

		client.Subscribe(channelName)

		channelColor := color.New(color.FgRed)
		eventColor := color.New(color.FgBlue)

		fmt.Printf("Subscribed to channel: %s\n", channelName)  // FIXME should be when we know it's successful

		client.BindGlobal(func(channelName string, eventName string, data interface{}) {
			fmt.Printf("Event: ")
			channelColor.Printf("channel=%s ", channelName)
			eventColor.Printf("event=%s ", eventName)
			fmt.Printf("message=%v\n", data)
		})

		fmt.Printf("Successfully subscribed to channel '")
		channelColor.Printf(channelName)
		fmt.Printf("'.\n")

		time.Sleep(time.Hour)
	},
}

func init() {
	Subscribe.PersistentFlags().StringVar(&appId, "app-id", "", "Pusher App ID")
	Subscribe.PersistentFlags().StringVar(&channelName, "channel", "", "Channel name")
}
