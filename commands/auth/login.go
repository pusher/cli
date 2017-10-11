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
		if len(apikey) < 1 {
			fmt.Println("Please first visit your Account settings page, and generate an API key.")
			os.Exit(1)
			return
		}
		fmt.Println("Got your API key!")
		conf.Token = apikey
		config.Store()

	},
}
