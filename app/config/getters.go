package config

import (
	"iac/utils/exitcodes"
	"log/slog"
	"os"

	"github.com/spf13/viper"
)

var cachedConfig Config

func GetConfig() Config {
	config := GetConfigUnvalidated()
	slog.Debug("Config loaded:", "config", config)

	err := validateConfig(&config)
	if err != nil {
		panic(err)
	}
	return config
}

func GetConfigUnvalidated() Config {
	if (cachedConfig != Config{}) {
		return cachedConfig
	}

	err := viper.ReadInConfig()
	if err != nil {
		slog.Error("Error reading config file", "error", err)
		os.Exit(exitcodes.ConfigFileNotFound)
	}

	err = viper.Unmarshal(&cachedConfig)
	if err != nil {
		panic(err)
	}

	return cachedConfig
}
