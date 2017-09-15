// Exposes all the needed configs for statsd-sink
package config

import (
	"os"

	"encoding/json"
	"io/ioutil"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

const (
	baseDirectory       = "~/.pusher"
	permissionOwnerOnly = 0200
)

type Config struct {
	Email    string `mapstructure:"email"`
	Password string `mapstructure:"password"`
}

var conf *Config

func readConfig(conf *Config) *Config {
	if _, err := os.Stat(baseDirectory); os.IsNotExist(err) {
		os.Mkdir(baseDirectory, permissionOwnerOnly)
	}

	viper.AddConfigPath(baseDirectory)
	viper.SetConfigName("config")

	viper.SetConfigType("json")

	if err := viper.ReadInConfig(); err != nil {
		logrus.WithError(err).Warnln("Failed to read config.")
	}

	if err := viper.Unmarshal(&conf); err != nil {
		logrus.WithError(err).Warnln("Failed to unmarshal config.")
	}

	return conf
}


func Get() *Config {
	if conf == nil {
		conf = readConfig(conf)
	}

	return conf
}

func Store() error {
	if _, err := os.Stat(baseDirectory); os.IsNotExist(err) {
		os.Mkdir(baseDirectory, permissionOwnerOnly)
	}

	confJson, err := json.Marshal(conf)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(baseDirectory+"/config.json", confJson, permissionOwnerOnly)
}

// IsSet checks if [configVariableName] has been set in the read config file.
func IsSet(configVariableName string) bool {
	return viper.IsSet(configVariableName)
}
