package api

import (
	"fmt"
	"os"

	"github.com/theherk/viper"
)

//isAPIKeyValid returns true if the stored API key is valid.
func isAPIKeyValid() bool {
	if viper.GetString("token") != "" {
		_, err := GetAllApps()
		if err == nil {
			return true
		}
	}
	return false
}

func validateKeyOrDie() {
	if !isAPIKeyValid() {
		fmt.Println("Your API key isn't valid. Add one with the `login` command.")
		os.Exit(1)
	}
}
