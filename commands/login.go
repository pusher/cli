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
		fmt.Println("What is your password? ")
	},
}
