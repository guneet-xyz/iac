package config

import (
	"github.com/spf13/viper"
)

var ConfigDirName = "~/.config/iac"
var ConfigFileNameWithoutExtension = "config"
var ConfigFileName = ConfigFileNameWithoutExtension + ".yml"
var ConfigFilePath string

func init() {
	err := setupConfigIfItDoesNotExist()
	if err != nil {
		panic(err)
	}

	viper.SetConfigName(ConfigFileNameWithoutExtension)
	viper.AddConfigPath(ConfigDirName)
}
