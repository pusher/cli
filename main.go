package main

import (
	"github.com/pusher/pusher-http-go"
	"github.com/urfave/cli"
	"os"
	"encoding/json"
	"io/ioutil"
	"log"
	"os/user"
	"path"
)

type appConfig struct {
	Key string `json:"key"`
	Secret string `json:"secret"`
	Cluster string `json:"cluster"`
}

type config struct {
	Apps map[string]appConfig `json:"apps"`
}

func main() {
	usr, err := user.Current()
	if err != nil {
		log.Fatalf("Could not get current user: %s", err.Error())
	}
	configFile, err := ioutil.ReadFile(path.Join(usr.HomeDir, ".pusher.json"))
	if err != nil {
		log.Fatalf("Could not read config: %s", err.Error())
	}
	var config config
	err = json.Unmarshal(configFile,&config)
	if err != nil {
		log.Fatalf("Could not parse config as JSON: %s", err.Error())
	}

	app := cli.NewApp()
	app.Name = "pusher"
	app.HelpName = "pusher"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name: "app-id",
			Usage: "Immutable ID for a Pusher app; find it on https://dashboard.pusher.com",
		},
	}
	app.Commands = []cli.Command{
		{
			Name: "trigger",
			Usage: "Trigger an event on a channel",
			ArgsUsage: "<channel-name> <event-name> <data>",
			Action: func (c *cli.Context) error {
				if !c.IsSet("app-id") {
					log.Fatalf("Please identify a Pusher app with --app-id")
				}
				appId := c.String("app-id")
				if appId == "" {
					log.Fatalf("--app-id cannot be empty string")
				}

				appConfig, appConfigExists := config.Apps[appId]
				if !appConfigExists {
					log.Fatalf("You have not set config for the app with id \"%s\". Please run: pusher add-app-config")
				}


				pusherApiClient := pusher.Client{
					AppId: appId,
					Key: appConfig.Key,
					Secret: appConfig.Secret,
					Cluster: appConfig.Cluster,
				}

				channelName := c.Args().Get(0)
				eventName := c.Args().Get(1)
				data := c.Args().Get(2)

				pusherApiClient.Trigger(channelName, eventName, data)

				return nil
			},
		},
	}

	app.Run(os.Args)


}