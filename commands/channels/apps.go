package channels

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"

	"github.com/olekukonko/tablewriter"
	"github.com/pusher/cli/api"
	"github.com/pusher/cli/commands"
	"github.com/pusher/cli/commands/auth"
	"github.com/spf13/cobra"
)

// Apps gets and displays a list of apps.
var Apps = &cobra.Command{
	Use:   "list",
	Short: "Get the list of all apps",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		if !auth.APIKeyValid() {
			fmt.Println("Your API key isn't valid. Add one with the `login` command.")
			os.Exit(1)
			return
		}

		apps, err := api.GetAllApps()
		if err != nil {
			fmt.Println("Failed to retrieve the list of apps.")
			os.Exit(1)
			return
		}

		if commands.OutputAsJSON {
			appsJSONBytes, _ := json.MarshalIndent(apps, "", "	")
			fmt.Println(string(appsJSONBytes))
		} else {
			table := tablewriter.NewWriter(os.Stdout)
			table.SetHeader([]string{"App ID", "App Name", "Cluster"})
			for _, app := range apps {
				table.Append([]string{strconv.Itoa(app.ID), app.Name, app.Cluster})
			}
			table.Render()
		}
	},
}

func init() {
	Apps.PersistentFlags().BoolVar(&commands.OutputAsJSON, "json", false, "")
}
