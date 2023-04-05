package channels

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/pusher/cli/api"
	"github.com/pusher/cli/commands"
	"github.com/pusher/pusher-http-go"
	"github.com/spf13/cobra"
	"github.com/theherk/viper"
)

var ChannelInfo = &cobra.Command{
	Use:   "channel-info",
	Short: "Get info about a specific channel",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {

		if commands.AppID == "" {
			fmt.Fprintf(os.Stderr, "Please supply --app-id\n")
			os.Exit(1)
			return
		}

		if commands.ChannelName == "" {
			fmt.Fprintf(os.Stderr, "Please supply --channel\n")
			os.Exit(1)
			return
		}

		p := api.NewPusherApi()
		app, err := p.GetApp(commands.AppID)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Could not get the app: %s\n", err.Error())
			os.Exit(1)
			return
		}

		token, err := p.GetToken(commands.AppID)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Could not get app token: %s\n", err.Error())
			os.Exit(1)
			return
		}

		client := pusher.Client{
			AppID:   commands.AppID,
			Key:     token.Key,
			Secret:  token.Secret,
			Cluster: app.Cluster,
			Host:    viper.GetString("apihost"),
		}

		infoValues := []string{}

		if commands.FetchUserCount {
			infoValues = append(infoValues, "user_count")
		}
		if commands.FetchSubscriptionCount {
			infoValues = append(infoValues, "subscription_count")
		}

		additionalQueries := map[string]string{}

		if len(infoValues) > 0 {
			infoString := infoValues[0]
			for i := 1; i < len(infoValues); i++ {
				infoString += "," + infoValues[i]
			}
			additionalQueries["info"] = infoString
		}

		channel, err := client.Channel(commands.ChannelName, additionalQueries)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Could not get channel: %s\n", err.Error())
			return
		}

		channelJSONBytes, _ := json.MarshalIndent(channel, "", "	")
		fmt.Println(string(channelJSONBytes))
	},
}

func init() {
	ChannelInfo.PersistentFlags().StringVar(&commands.AppID, "app-id", "", "Channels App ID")
	ChannelInfo.PersistentFlags().StringVar(&commands.ChannelName, "channel", "", "Channel name")
	ChannelInfo.PersistentFlags().BoolVar(&commands.FetchUserCount, "fetch-user-count", false, "Fetch user count for the presence channel")
	ChannelInfo.PersistentFlags().BoolVar(&commands.FetchSubscriptionCount, "fetch-subscription-count", false, "Fetch subscription count for the channel")
}
