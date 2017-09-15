package commands

import (
	"github.com/pusher-community/pusher-websocket-go"
	"github.com/spf13/cobra"
	"fmt"
	"os"
	"github.com/pusher/pusher-cli/api"
	"time"
)

var Subscribe = &cobra.Command{
	Use:   "subscribe [channel]",
	Short: "Subscribe to a Pusher channel",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
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

		channelName := args[0]
		client.Subscribe(channelName)

		client.BindGlobal(func (channelName string, eventName string, data interface{}) {
			fmt.Printf("Event received. channel=%s event=%s message=%v\n", channelName, eventName, data)
		})

		time.Sleep(time.Hour)
	},
}

func init() {
	Subscribe.PersistentFlags().StringVar(&appId, "app-id", "", "Pusher App ID")
}
