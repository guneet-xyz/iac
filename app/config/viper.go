package config

import (
	"github.com/spf13/viper"
	"iac/utils/fs"
)

var ConfigDirName = "~/.config/iac"
var ConfigFileNameWithoutExtension = "config"
var ConfigFileName = ConfigFileNameWithoutExtension + ".yml"
var ConfigFilePath string

func init() {
	ConfigDirName, err := fs.AbsPath(ConfigDirName)
	if err != nil {
		panic(err)
	}
	viper.SetConfigName(ConfigFileNameWithoutExtension)
	viper.AddConfigPath(ConfigDirName)
}
