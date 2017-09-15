package commands

import (
	"github.com/pusher-community/pusher-websocket-go"
	"github.com/spf13/cobra"
	"fmt"
	"os"
	"github.com/pusher/pusher-cli/api"
	"time"
	"github.com/fatih/color"
)


var Subscribe = &cobra.Command{
	Use:   "subscribe",
	Short: "Subscribe to a Pusher channel",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {

		if appId == "" {
			fmt.Fprintf(os.Stderr,"Please supply --app-id\n")
			return
		}

		if channelName == "" {
			fmt.Fprintf(os.Stderr,"Please supply --channel\n")
			return
		}

		token, err := api.GetToken(appId)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Could not get token: %s\n", err.Error())
			return
		}

		pusher.New("foo")
		client := pusher.NewWithConfig(pusher.ClientConfig{
			Scheme: "wss",
			Host: "ws-test1.staging.pusher.com",
			Port: "443",
			Key: token.Key,
			Secret: token.Secret,
		})

		client.Subscribe(channelName)

		channelColor := color.New(color.FgCyan).Add(color.Underline)

		client.BindGlobal(func (channelName string, eventName string, data interface{}) {
			fmt.Printf("Event: ")
			channelColor.Printf("Testing")
			fmt.Printf("Event: channel=%s event=%s message=%v\n", channelName, eventName, data)
		})

		time.Sleep(time.Hour)
	},
}

func init() {
	Subscribe.PersistentFlags().StringVar(&appId, "app-id", "", "Pusher App ID")
	Subscribe.PersistentFlags().StringVar(&channelName, "channel", "", "Channel name")
}
