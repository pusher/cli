package auth

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/theherk/viper"
)

// Logout removes the users API key from the machine.
var Logout = &cobra.Command{
	Use:   "logout",
	Short: "Remove Pusher account credentials from this computer",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		viper.Set("token", "")
		err := viper.WriteConfig()
		if err != nil {
			panic("Could not write config: " + err.Error())
		}
		fmt.Println("Removed Pusher account credentials.")
	},
}
