package channels

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/pusher/cli/api"
	"github.com/pusher/cli/commands"
	"github.com/spf13/cobra"
)

//GenerateClient generates a client that can subscribe to channels on an app.
var GenerateClient = &cobra.Command{
	Use:   "client",
	Short: "Generate a client for your Channels app",
	Args:  cobra.NoArgs,
}

//GenerateWeb generates a web client that can subscribe to channels on an app.
var GenerateWeb = &cobra.Command{
	Use:   "web",
	Short: "Generate a web client for your Channels app",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {

		if commands.AppID == "" {
			fmt.Fprintf(os.Stderr, "Please supply --app-id\n")
			return
		}

		token, err := api.GetToken(commands.AppID)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Could not get app token: %s\n", err.Error())
			return
		}

		html := `
			<!DOCTYPE html>
			<head>
			<title>Pusher Test</title>
			<script src="https://js.pusher.com/4.1/pusher.min.js"></script>
			<script>
			// Enable pusher logging - don't include this in production
			Pusher.logToConsole = true;
			
			var pusher = new Pusher('` + token.Key + `', {
      	wsHost: 'ws-test1.staging.pusher.com',
      	httpHost: 'sockjs-test1.staging.pusher.com',
				encrypted: true
			});
			
			var channel = pusher.subscribe('my-channel');
    	channel.bind('my-event', function(data) {
				alert(data.message);
			});
			</script>
			</head>
		`
		err = ioutil.WriteFile("index.html", []byte(html), 0644)
		if err != nil {
			fmt.Printf("Could not write file: %s\n", err.Error())
		} else {
			fmt.Printf("Written file: index.html\n")
		}
	},
}

func init() {
	GenerateClient.PersistentFlags().StringVar(&commands.AppID, "app-id", "", "Channels App ID")
	GenerateClient.AddCommand(GenerateWeb)
}
