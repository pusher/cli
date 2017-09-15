package commands

import (
	"fmt"
	"github.com/spf13/cobra"
)

var Login = &cobra.Command{
	Use:   "login",
	Short: "Enter and store Pusher account credentials",
	Args: cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("What is your email address? ")
		var email string
		fmt.Scanln(&email)
		fmt.Println("What is your password? ")
		var password string
		fmt.Scanln(&password)
		fmt.Printf("Email: %s, Password: %s\n", email, password)
	},
}
