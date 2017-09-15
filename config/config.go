// Exposes all the needed configs for statsd-sink
package config

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"os/user"

	"github.com/spf13/viper"
)

const (
	baseDirectory           = "/.pusher"
	defaultDirPermission    = 0755
	defaultConfigPermission = 0600
)

type Config struct {
	Email    string `mapstructure:"email"`
	Password string `mapstructure:"password"`
}

var conf *Config

func getUserHomeDir() string {
	usr, err := user.Current()
	if err != nil {
		return "~" // best effort
	}

	return usr.HomeDir
}

func readConfig() *Config {
	if _, err := os.Stat(baseDirectory); os.IsNotExist(err) {
		os.Mkdir(baseDirectory, defaultDirPermission)
	}

	viper.AddConfigPath(getUserHomeDir() + baseDirectory)
	viper.SetConfigName("config")

	viper.SetConfigType("json")

	c := &Config{}
	if err := viper.ReadInConfig(); err != nil {
		//logrus.WithError(err).Warnln("Failed to read config.")
		return c
	}

	if err := viper.Unmarshal(&c); err != nil {
		//logrus.WithError(err).Warnln("Failed to unmarshal config.")
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
		os.Mkdir(getUserHomeDir()+baseDirectory, defaultDirPermission)
	}

	confJson, err := json.Marshal(conf)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(getUserHomeDir()+baseDirectory+"/config.json", confJson, defaultConfigPermission)
}

// IsSet checks if [configVariableName] has been set in the read config file.
func IsSet(configVariableName string) bool {
	return viper.IsSet(configVariableName)
}
