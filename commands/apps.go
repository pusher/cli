package commands

import (
	"fmt"

	"os"

	"github.com/pusher/pusher-cli/api"
	"github.com/pusher/pusher-cli/config"
	"github.com/spf13/cobra"
)

var Apps = &cobra.Command{
	Use:   "apps",
	Short: "Get the list of all apps",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		if config.Get().Email == "" || config.Get().Password == "" {
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

		for _, app := range apps {
			fmt.Printf("%+v\n", app)
		}
	},
}
