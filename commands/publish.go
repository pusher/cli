package commands

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/pusher/pusher-http-go"
	"os"
)

type AppToken struct {
	Key string
	Secret string
}

type App struct {
	AppId string
	Cluster string
	Tokens []AppToken
}

func GetApp(appId string) *App {
	return &App{AppId: appId, Cluster: "mt1", Tokens:[]AppToken{AppToken{"foo", "bar"}}}
}

var appId string
var channelName string
var eventName string
var message string

var Publish = &cobra.Command{
	Use:   "publish [channel]",
	Short: "Trigger an event on a Pusher app",
	Args: cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		app := GetApp(appId)
		token := app.Tokens[0]
		client := pusher.Client{
			AppId: appId,
			Key: token.Key,
			Secret: token.Secret,
			Cluster: app.Cluster,
		}

		_, err := client.Trigger(channelName, eventName, message)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Could not trigger: %s\n", err.Error())
			return
		}

		fmt.Printf("Triggered event: %s %s %s %s\n", appId, channelName, eventName, message)
	},
}

func init() {
	Publish.PersistentFlags().StringVar(&appId, "app-id", "", "Pusher App ID")
	Publish.PersistentFlags().StringVar(&channelName, "channel", "", "Channel name")
	Publish.PersistentFlags().StringVar(&channelName, "event", "", "Event name")
	Publish.PersistentFlags().StringVar(&channelName, "message", "", "Message")
}

