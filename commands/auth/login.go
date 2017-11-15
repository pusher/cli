package auth

import (
	"fmt"
	"os"

	"golang.org/x/crypto/ssh/terminal"

	"github.com/pusher/cli/api"
	"github.com/spf13/cobra"
	"github.com/theherk/viper"
)

// Login allows users to log in using an API token.
var Login = &cobra.Command{
	Use:   "login",
	Short: "Enter and store Pusher account credentials",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		if APIKeyValid() {
			fmt.Println("Your API key is valid. If you'd like to use a different API key, use `logout` first.")
			os.Exit(1)
		}
		fmt.Println("What is your email address?")
		var email string
		fmt.Scanln(&email)

		fmt.Println("What is your password?")
		passwordBytes, _ := terminal.ReadPassword(0)
		password := string(passwordBytes)

		// check if the user/pass can get an API key
		apikey, err := api.GetAPIKey(email, password)
		if err != nil {
			panic("Could not get API key: " + err.Error())
		}
		if apikey == "" {
			fmt.Println("There is No API key associated with those account details. Make sure you've set up your API key in the Admin Dashboard, and that are your details are correct.")
			return
		}
		fmt.Println("Got your API key!")
		viper.Set("token", apikey)
		fmt.Println(apikey)
		err = viper.WriteConfig()
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
