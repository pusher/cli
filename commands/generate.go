package commands

import (
	"fmt"
	"github.com/spf13/cobra"
	"io/ioutil"
	"os"
	"github.com/pusher/pusher-cli/api"
)

var Generate = &cobra.Command{
	Use:   "generate-client",
	Short: "Generate a Pusher client for your Pusher app",
	Args:  cobra.NoArgs,
}

var GenerateWeb = &cobra.Command{
	Use:   "web",
	Short: "Generate a web client for your Pusher app",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {

		if appId == "" {
			fmt.Fprintf(os.Stderr,"Please supply --app-id\n")
			return
		}

		token, err := api.GetToken(appId)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Could not get token: %s\n", err.Error())
			return
		}

		html :=
`<!DOCTYPE html>
<head>
  <title>Pusher Test</title>
  <script src="https://js.pusher.com/4.1/pusher.min.js"></script>
  <script>

    // Enable pusher logging - don't include this in production
    Pusher.logToConsole = true;

    var pusher = new Pusher('`+token.Key+`', {
      wsHost: 'ws-test1.staging.pusher.com',
      httpHost: 'sockjs-test1.staging.pusher.com',
      encrypted: true
    });

    var channel = pusher.subscribe('my-channel');
    channel.bind('my-event', function(data) {
      alert(data.message);
    });
  </script>
</head>`
		err = ioutil.WriteFile("index.htm", []byte(html), 0644)
		if err != nil {
			fmt.Printf("Could not write file: %s\n", err.Error())
		} else {
			fmt.Printf("Written file: index.htm\n")
		}
	},
}

func init() {
	Generate.PersistentFlags().StringVar(&appId, "app-id", "", "Pusher App ID")
	Generate.AddCommand(GenerateWeb)
}