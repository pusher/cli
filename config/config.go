package config

import (
	"fmt"
	"os"
	"os/user"
	"path"

	"io/ioutil"

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

func getConfigDir() string {
	return path.Join(getUserHomeDir(), ".config2")
}

func getConfigPath() string {
	return path.Join(getConfigDir(), "pusher.json")
}

//Init sets the config files location and attempts to read it in.
func Init() {
	if _, err := os.Stat(getConfigDir()); os.IsNotExist(err) {
		err = os.Mkdir(getConfigDir(), os.ModeDir|0755)
		if err != nil {
			panic(err)
		}
	}

	if _, err := os.Stat(getConfigPath()); os.IsNotExist(err) {
		err = ioutil.WriteFile(getConfigPath(), []byte("{}"), 0600)
		if err != nil {
			panic(err)
		}
	}

	viper.SetConfigFile(getConfigPath())
	viper.SetDefault("endpoint", "https://cli.pusher.com")
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
}
