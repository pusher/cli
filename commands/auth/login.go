package auth

import (
	"fmt"
	"os"

	"github.com/pusher/pusher-cli/api"
	"github.com/pusher/pusher-cli/config"
	"github.com/spf13/cobra"
)

// Login allows users to log in using an API token.
var Login = &cobra.Command{
	Use:   "login",
	Short: "Enter and store Pusher account credentials",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		conf := config.Get()

		if APIKeyValid() {
			fmt.Println("Your API key is valid. If you'd like to use a different API key, use `logout` first.")
			os.Exit(1)
			return
		}

		fmt.Println("What is your Pusher API Key? You can get this by visiting your account settings page (http://localhost:3001/accounts/edit)")
		fmt.Scanln(&conf.Token)
		// check if the credentials are valid
		_, err := api.GetAllApps()
		if err != nil {
			fmt.Println("Invalid credentials.")
			os.Exit(1)
		} else {
			config.Store()
			fmt.Println("Successfully logged in.")
		}
	},
}
