package config

import (
	"fmt"
	"os"
	"os/user"

	"github.com/theherk/viper"
)

const (
	baseDirectory   = "/.config/"
	fileName        = "pusherconfig.json"
	defaultEndpoint = "http://localhost:3000"
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

//Init sets the config files location and attempts to read it in.
func Init() {
	viper.SetConfigFile(getUserHomeDir() + baseDirectory + fileName)
	viper.SetDefault("endpoint", defaultEndpoint)
	viper.ReadInConfig()
}

// Delete deletes the config file.
func Delete() error {
	return os.Remove(getUserHomeDir() + baseDirectory + fileName)
}
