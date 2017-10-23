package auth

import (
	"fmt"
	"os"

	"golang.org/x/crypto/ssh/terminal"

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
		var e string
		fmt.Println("What is your email address? ")

		fmt.Scanln(&e)

		fmt.Println("What is your password? ")
		passwordBytes, _ := terminal.ReadPassword(0)
		p := string(passwordBytes)
		// check if the user/pass can get an API key
		apikey, err := api.GetAPIKey(e, p)
		if err != nil {
			fmt.Println("Error getting API key")
			os.Exit(1)
			return
		}
		if apikey == "" {
			fmt.Println("There is No API key associated with those account details. Make sure you've set up your API key in the Admin Dashboard, and that are your details are correct.")
			return
		}
		fmt.Println("Got your API key!")
		conf.Token = apikey
		config.Store()

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
