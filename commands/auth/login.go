package auth

import (
	"fmt"
	"os"

	"github.com/pusher/cli/api"
	"github.com/spf13/cobra"
	"github.com/theherk/viper"
)

// Login allows users to log in using an API token.
var Login = &cobra.Command{
	Use:   "login",
	Short: "Enter and store Pusher API key",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		if APIKeyValid() {
			fmt.Println("Your current API key is valid. If you'd like to use a different API key, use `logout` first.")
			os.Exit(1)
		}

		fmt.Println("What is your API key? (Find this at https://dashboard.pusher.com/accounts/edit)")
		var apikey string
		fmt.Scanln(&apikey)
		viper.Set("token", apikey)

		if !APIKeyValid() {
			fmt.Println("That API key is not valid!")
			os.Exit(1)
		}

		err := viper.WriteConfig()
		if err != nil {
			panic("Could not write config: " + err.Error())
		}
	},
}

//APIKeyValid returns true if the stored API key is valid.
func APIKeyValid() bool {
	if viper.GetString("token") != "" {
		_, err := api.GetAllApps()
		if err == nil {
			return true
		}
	}
	return false
}
