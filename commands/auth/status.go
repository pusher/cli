package auth

import (
	"fmt"

	"github.com/pusher/pusher-cli/api"
	"github.com/pusher/pusher-cli/config"
	"github.com/spf13/cobra"
)

// CheckAuth checks the stored API key for validity
var Status = &cobra.Command{
	Use:   "status",
	Short: "Checks the stored API key for validity.",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		if APIKeyValid() {
			fmt.Println("API Key is valid.")
		} else {
			fmt.Println("API Key is NOT valid.")
		}
	},
}

//APIKeyValid returns true if the stored API key is valid.
func APIKeyValid() bool {
	conf := config.Get()
	if conf.Token != "" {
		_, err := api.GetAllApps()
		if err == nil {
			return true
		}
	}
	return false
}
