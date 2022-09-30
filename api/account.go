package api

import (
	"fmt"
	"os"

	"github.com/theherk/viper"
)

//isAPIKeyValid returns true if the stored API key is valid.
func (p *PusherApi) isAPIKeyValid() bool {
	if viper.GetString("token") != "" {
		_, err := p.GetAllApps()
		if err == nil {
			return true
		}
	}
	return false
}

func (p *PusherApi) validateKeyOrDie() {
	if !p.isAPIKeyValid() {
		fmt.Println("Your API key isn't valid. Add one with the `login` command.")
		os.Exit(1)
	}
}
