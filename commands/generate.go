package commands

import (
	"fmt"
	"github.com/spf13/cobra"
)

var Generate = &cobra.Command{
	Use:   "generate",
	Short: "Generate a starter project which uses a Pusher app",
	Args:  cobra.NoArgs,
}

var GenerateWeb = &cobra.Command{
	Use:   "web",
	Short: "Generate a starter web project which uses a Pusher app",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("TODO Generating a web project")
	},
}

func init() {
	Generate.PersistentFlags().StringVar(&appId, "app-id", "", "Pusher App ID")
	Generate.AddCommand(GenerateWeb)
}