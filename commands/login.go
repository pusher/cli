package commands

import (
	"fmt"
	"os"

	"github.com/pusher/pusher-cli/api"
	"github.com/pusher/pusher-cli/config"
	"github.com/spf13/cobra"
	"golang.org/x/crypto/ssh/terminal"
)

var Login = &cobra.Command{
	Use:   "login",
	Short: "Enter and store Pusher account credentials",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		conf := config.Get()

		if config.Get().Email != "" && config.Get().Password != "" {
			fmt.Printf("Already logged in as '%s'.\n", config.Get().Email)
			os.Exit(1)
			return
		}

		fmt.Println("What is your email address? ")
		fmt.Scanln(&conf.Email)
		fmt.Println("What is your password? ")
		passwordBytes, _ := terminal.ReadPassword(0)

		conf.Password = string(passwordBytes)

		// check if the credentials are valid
		_, err := api.GetAllApps()
		if err != nil {
			fmt.Println("Invalid credentials.")
			os.Exit(1)
		} else {
			config.Store()
			fmt.Println("Succesfully logged in.")
		}
	},
}
