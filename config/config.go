package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/user"

	"path"

	"github.com/spf13/viper"
)

const (
	defaultDirPermission    = 0755
	defaultConfigPermission = 0600
)

type Config struct {
	Token string `mapstructure:"token"`
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

func configDirPath() string {
	return path.Join(getUserHomeDir(), ".config/pusher-cli")
}

func configJsonPath() string {
	return path.Join(configDirPath(), "config.json")
}

func ensureConfigDirExists() {
	if _, err := os.Stat(configDirPath()); os.IsNotExist(err) {
		err := os.MkdirAll(configDirPath(), defaultDirPermission)
		if err != nil {
			panic("Could not create config directory: " + err.Error())
		}
	}
}

func readConfig() *Config {
	ensureConfigDirExists()

	viper.AddConfigPath(configDirPath())
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
	ensureConfigDirExists()

	confJson, err := json.Marshal(conf)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(configJsonPath(), confJson, defaultConfigPermission)
}

func Delete() error {
	return os.Remove(configJsonPath())
}

// IsSet checks if [configVariableName] has been set in the read config file.
func IsSet(configVariableName string) bool {
	return viper.IsSet(configVariableName)
}
