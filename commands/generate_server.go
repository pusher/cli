package commands

import (
	"fmt"
	"github.com/spf13/cobra"
	"io/ioutil"
	"os"
	"github.com/pusher/pusher-cli/api"
)

var GenerateServer = &cobra.Command{
	Use:   "generate-server",
	Short: "Generate a Pusher server for your Pusher app",
	Args:  cobra.NoArgs,
}

var GeneratePhp = &cobra.Command{
	Use:   "php",
	Short: "Generate a PHP server for your Pusher app",
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

		php :=
`<?php
  require __DIR__ . '/vendor/autoload.php';

  $options = array(
    'host' => 'api-test1.staging.pusher.com',
    'encrypted' => true
  );
  $pusher = new Pusher(
    '`+token.Key+`',
    '`+token.Secret+`',
    '`+appId+`',
    $options
  );

  $data['message'] = 'hello world';
  $pusher->trigger('my-channel', 'my-event', $data);
?>`
		err = ioutil.WriteFile("main.php", []byte(php), 0644)
		if err != nil {
			fmt.Printf("Could not write file: %s\n", err.Error())
		} else {
			fmt.Printf("Written file: main.php.\nPlease run: composer require pusher/pusher-php-server\n")
		}
	},
}

var GeneratePython = &cobra.Command{
	Use:   "python",
	Short: "Generate a Python server for your Pusher app",
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

		python :=
    `import pusher

pusher_client = pusher.Pusher(
  app_id='`+appId+`',
  key='`+token.Key+`',
  secret='`+token.Secret+`',
  host='api-test1.staging.pusher.com',
  ssl=True
)

pusher_client.trigger('my-channel', 'my-event', {'message': 'hello world'})`
		err = ioutil.WriteFile("server.py", []byte(python), 0644)
		if err != nil {
			fmt.Printf("Could not write file: %s\n", err.Error())
		} else {
			fmt.Printf("Written file: server.py.\nPlease run: pip install pusher\n")
		}
	},
}

func init() {
	GenerateServer.PersistentFlags().StringVar(&appId, "app-id", "", "Pusher App ID")
	GenerateServer.AddCommand(GeneratePhp)
	GenerateServer.AddCommand(GeneratePython)
}