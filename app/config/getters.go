package config

import (
	"iac/utils/exitcodes"
	"log/slog"
	"os"

	"github.com/spf13/viper"
)

var cachedConfig Config

func GetConfig() (Config, error) {
	config, err := GetConfigUnvalidated()
	if err != nil {
		slog.Error("Error getting config", "error", err)
		return Config{}, err
	}
	slog.Debug("Config loaded:", "config", config)

	err = validateConfig(&config)
	if err != nil {
		slog.Error("Config validation failed", "error", err)
		return Config{}, err
	}

	return config, nil
}

func GetConfigUnvalidated() (Config, error) {
	if (cachedConfig != Config{}) {
		return cachedConfig, nil
	}

	err := viper.ReadInConfig()
	if err != nil {
		slog.Error("Error reading config file", "error", err)
		os.Exit(exitcodes.ConfigFileNotFound)
	}

	err = viper.Unmarshal(&cachedConfig)
	if err != nil {
		slog.Error("Error unmarshaling config file", "error", err)
		return Config{}, err
	}

	return cachedConfig, nil
}
