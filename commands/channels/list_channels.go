package channels

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/pusher/cli/api"
	"github.com/pusher/cli/commands"
	"github.com/pusher/pusher-http-go"
	"github.com/spf13/cobra"
)

var ListChannels = &cobra.Command{
	Use:   "list-channels",
	Short: "List all channels that have at least one subscriber",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {

		if commands.AppID == "" {
			fmt.Fprintf(os.Stderr, "Please supply --app-id\n")
			os.Exit(1)
			return
		}

		app, err := api.GetApp(commands.AppID)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Could not get app the app: %s\n", err.Error())
			os.Exit(1)
			return
		}

		token, err := api.GetToken(commands.AppID)
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
		}

		opts := map[string]string{}

		if commands.FilterByPrefix != "" {
			opts["filter_by_prefix"] = commands.FilterByPrefix
		}

		channelsList, err := client.Channels(opts)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Could not get channel list: %s\n", err.Error())
			return
		}

		channelsListJSONBytes, _ := json.MarshalIndent(channelsList.Channels, "", "	")
		fmt.Println(string(channelsListJSONBytes))
	},
}

func init() {
	ListChannels.PersistentFlags().StringVar(&commands.AppID, "app-id", "", "Channels App ID")
	ListChannels.PersistentFlags().StringVar(&commands.FilterByPrefix, "filter-by-prefix", "", "A channel name prefix, e.g. 'presence-'")
}
