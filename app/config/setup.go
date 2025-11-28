package config

import (
	_ "embed"
	"iac/utils/fs"
	"log/slog"
	"path/filepath"
)

//go:embed default_config.toml
var DefaultConfig string

func setupConfigIfItDoesNotExist() error {
	var err error

	configDirName := filepath.Dir(ConfigFilePath)
	_, err = fs.MkdirIfNotExists(configDirName)
	if err != nil {
		panic(err)
	}

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
