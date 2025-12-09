package config

import (
	"log/slog"
	"os"
	"strings"

	"iac/utils/exitcodes"
	"iac/utils/fs"

	"github.com/spf13/viper"
)

var ConfigFilePath = "~/.config/iac/config.toml"

func init() {
	var err error
	envConfigPath := os.Getenv("IAC_CONFIG")

	if envConfigPath != "" {
		if !strings.HasSuffix(strings.ToLower(envConfigPath), ".toml") {
			slog.Error("Environment variable IAC_CONFIG must point to a .toml file")
			os.Exit(exitcodes.BadEnvironmentVariables)
		}
		ConfigFilePath = envConfigPath
	}

	ConfigFilePath, err = fs.AbsPath(ConfigFilePath)
	if err != nil {
		slog.Error("Failed to resolve config file path", "error", err)
		os.Exit(exitcodes.CouldNotResolveConfigFilePath)
	}

	err = setupConfigIfItDoesNotExist()
	if err != nil {
		slog.Error("Failed to setup config file", "error", err)
		os.Exit(exitcodes.CouldNotSetupConfigFile)
	}

	viper.SetConfigFile(ConfigFilePath)
}
