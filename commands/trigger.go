package commands

import (
	"fmt"
	"github.com/pusher/pusher-http-go"
	"github.com/spf13/cobra"
	"os"
	"github.com/pusher/pusher-cli/api"
)

type AppToken struct {
	Key    string
	Secret string
}

type App struct {
	AppId   string
	Cluster string
	Tokens  []AppToken
}

func GetApp(appId string) *App {
	return &App{AppId: appId, Cluster: "mt1", Tokens: []AppToken{AppToken{"foo", "bar"}}}
}

var appId string
var channelName string
var eventName string
var message string

var Trigger = &cobra.Command{
	Use:   "trigger",
	Short: "Trigger an event on a Pusher app",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {

		token, err := api.GetToken(appId)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Could not get token: %s\n", err.Error())
			return
		}

		client := pusher.Client{
			AppId:   appId,
			Key:     token.Key,
			Secret:  token.Secret,
			Cluster: "eu", // app.Cluster,
		}

		fmt.Printf("Triggering event: %s %s %s %s\n", appId, channelName, eventName, message)

		_, err = client.Trigger(channelName, eventName, message)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Could not trigger: %s\n", err.Error())
			return
		}

		fmt.Printf("Triggered event: %s %s %s %s\n", appId, channelName, eventName, message)
	},
}

func init() {
	Trigger.PersistentFlags().StringVar(&appId, "app-id", "", "Pusher App ID")
	Trigger.PersistentFlags().StringVar(&channelName, "channel", "", "Channel name")
	Trigger.PersistentFlags().StringVar(&eventName, "event", "", "Event name")
	Trigger.PersistentFlags().StringVar(&message, "message", "", "Message")
}
