package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/user"

	"github.com/spf13/viper"
)

const (
	baseDirectory           = "/.config/pusher-cli"
	defaultDirPermission    = 0755
	defaultConfigPermission = 0600
)

type Config struct {
	Token        string `mapstructure:"token"`
	BaseEndpoint string `mapstructure:"baseEndpoint"`
}

var conf *Config
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

func readConfig() *Config {
	if _, err := os.Stat(baseDirectory); os.IsNotExist(err) {
		os.MkdirAll(baseDirectory, defaultDirPermission)
	}

	viper.AddConfigPath(getUserHomeDir() + baseDirectory)
	viper.SetConfigName("config")

	viper.SetConfigType("json")

	c := &Config{}
	if err := viper.ReadInConfig(); err != nil {
		return c
	}

	if err := viper.Unmarshal(&c); err != nil {
		return c
	}

	return c
}

func Get() *Config {
	if conf == nil {
		conf = readConfig()
	}

	return conf
}

func Store() error {
	if _, err := os.Stat(getUserHomeDir() + baseDirectory); os.IsNotExist(err) {
		os.MkdirAll(getUserHomeDir()+baseDirectory, defaultDirPermission)
	}

	confJson, err := json.Marshal(conf)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(getUserHomeDir()+baseDirectory+"/config.json", confJson, defaultConfigPermission)
}

func Delete() error {
	return os.Remove(getUserHomeDir() + baseDirectory + "/config.json")
}

// IsSet checks if [configVariableName] has been set in the read config file.
func IsSet(configVariableName string) bool {
	return viper.IsSet(configVariableName)
}
