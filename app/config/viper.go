package config

import (
	"os"
	"strings"

	"iac/utils/fs"

	"github.com/spf13/viper"
)

var ConfigFilePath = "~/.config/iac/config.toml"

func init() {
	var err error
	envConfigPath := os.Getenv("IAC_CONFIG")

	if envConfigPath != "" {
		if !strings.HasSuffix(strings.ToLower(envConfigPath), ".toml") {
			panic("IAC_CONFIG must be a .toml file")
		}
		ConfigFilePath = envConfigPath
	}

	ConfigFilePath, err = fs.AbsPath(ConfigFilePath)
	if err != nil {
		panic(err)
	}

	err = setupConfigIfItDoesNotExist()
	if err != nil {
		panic(err)
	}

	viper.SetConfigFile(ConfigFilePath)
}
