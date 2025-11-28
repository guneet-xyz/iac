package config

import (
	_ "embed"
	"iac/utils/exitcodes"
	"iac/utils/fs"
	"log/slog"
	"os"
	"path/filepath"
)

//go:embed default_config.toml
var DefaultConfig string

func setupConfigIfItDoesNotExist() error {
	var err error

	ConfigDirName, err = fs.MkdirIfNotExists(ConfigDirName)
	if err != nil {
		panic(err)
	}

	ConfigFilePath := filepath.Join(ConfigDirName, ConfigFileName)

	stat, err := fs.Stat(ConfigFilePath)
	if err != nil || stat == fs.StatResultStatError {
		return err
	}

	switch stat {
	case fs.StatResultNotExist:
		slog.Info("Config file does not exist, creating default config file", "path", ConfigFilePath)
		err = fs.WriteFileFromString(ConfigFilePath, DefaultConfig)
	}

	return nil
}

func init() {
	err := setupConfigIfItDoesNotExist()
	if err != nil {
		slog.Error("Failed to setup config file", "error", err)
		os.Exit(exitcodes.FailedToSetupConfigFile)
	}
}
