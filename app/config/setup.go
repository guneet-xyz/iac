package config

import (
	_ "embed"
	"iac/utils/exitcodes"
	"iac/utils/fs"
	"log/slog"
	"path/filepath"
)

//go:embed templates/user_config.toml
var DefaultUserConfig string

//go:embed templates/repository_config.toml
var DefaultRepoConfig string

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
		err = fs.WriteFileFromString(ConfigFilePath, DefaultUserConfig)
	}

	return err
}

func SetupRepoConfigIfNotExists(repoPath string) error {
	repoDirStat, err := fs.Stat(repoPath)
	if err != nil || repoDirStat == fs.StatResultStatError {
		return err
	}

	if repoDirStat == fs.StatResultNotExist {
		slog.Info("Repository directory does not exist, creating it", "path", repoPath)
		_, err = fs.MkdirIfNotExists(repoPath)
		if err != nil {
			return err
		}
	}

	repoConfigPath := filepath.Join(repoPath, "config.toml")

	stat, err := fs.Stat(repoConfigPath)
	if err != nil || stat == fs.StatResultStatError {
		return err
	}

	if stat == fs.StatResultNotExist {
		slog.Info("Repository config file does not exist, creating default config file", "path", repoConfigPath)
		if err = fs.WriteFileFromString(repoConfigPath, DefaultRepoConfig); err != nil {
			return err
		}
	}

	return nil
}

func SanityChecks() (int, error) {
	slog.Debug("Performing config sanity checks")

	err := setupConfigIfItDoesNotExist()
	if err != nil {
		slog.Error("Failed to setup config file", "error", err)
		return exitcodes.FailedToSetupConfigFile, err
	}

	slog.Debug("Config sanity checks passed")
	return 0, nil
}
