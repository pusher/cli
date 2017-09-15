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

var Publish = &cobra.Command{
	Use:   "publish [app-id] [channel] [event] [message]",
	Short: "Trigger an event on a Pusher app",
	Args: cobra.ExactArgs(4),
	Run: func(cmd *cobra.Command, args []string) {
		appId := args[0]
		channelName := args[1]
		eventName := args[2]
		message := args[3]

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
