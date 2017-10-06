package commands

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"

	"github.com/olekukonko/tablewriter"
	"github.com/pusher/pusher-cli/api"
	"github.com/pusher/pusher-cli/config"
	"github.com/spf13/cobra"
)

// Apps gets and displays a list of apps.
var Apps = &cobra.Command{
	Use:   "apps",
	Short: "Get the list of all apps",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		if config.Get().Token == "" {
			fmt.Println("Not logged in.")
			os.Exit(1)
			return
		}

		apps, err := api.GetAllApps()
		if err != nil {
			fmt.Println("Failed to retrieve the list of apps.")
			os.Exit(1)
			return
		}

		if outputAsJSON {
			appsJSONBytes, _ := json.Marshal(apps)
			fmt.Println(string(appsJSONBytes))
		} else {
			table := tablewriter.NewWriter(os.Stdout)
			table.SetHeader([]string{"App ID", "App Name", "Cluster"})
			for _, app := range apps {
				table.Append([]string{strconv.Itoa(app.Id), app.Name, app.Cluster})
			}
			table.Render()
		}
	},
}

func init() {
	Apps.PersistentFlags().BoolVar(&outputAsJSON, "json", false, "")
}
