package config

import (
	"fmt"
	"os"
	"os/user"
	"path"

	"github.com/theherk/viper"
)

var version = "master"

// GetVersion returns version of PusherCLI, set in ldflags.
func GetVersion() string {
	return version
}

func getUserHomeDir() string {
	usr, err := user.Current()
	if err != nil {
		fmt.Println("Can't get your home directory.")
		os.Exit(1)
	}

	return usr.HomeDir
}

func getConfigPath() string {
	return path.Join(getUserHomeDir(), ".config/pusher.json")
}

//Init sets the config files location and attempts to read it in.
func Init() {
	viper.SetConfigFile(getConfigPath())
	viper.SetDefault("endpoint", "https://cli.pusher.com")
	viper.ReadInConfig()
}

// Delete deletes the config file.
func Delete() error {
	return os.Remove(getConfigPath())
}
