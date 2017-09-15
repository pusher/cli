package commands

import (
	"fmt"

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

		fmt.Println("What is your email address? ")
		fmt.Scanln(&conf.Email)
		fmt.Println("What is your password? ")
		passwordBytes, _ := terminal.ReadPassword(0)

		conf.Password = string(passwordBytes)

		fmt.Println(config.Store())
	},
}
