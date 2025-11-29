package config

import (
	"iac/utils/exitcodes"
	"iac/utils/fs"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

var cachedUserConfig UserConfig
var cachedRepoConfig RepoConfig
var cachedRepoPath string

func getUserConfigUnvalidated() UserConfig {
	err := viper.ReadInConfig()
	if err != nil {
		slog.Error("Error reading config file", "error", err)
		os.Exit(exitcodes.ConfigFileNotFound)
	}

	var userConfig UserConfig
	err = viper.Unmarshal(&userConfig)
	if err != nil {
		panic(err)
	}

	return userConfig
}

func GetUserConfig() UserConfig {
	if (cachedUserConfig != UserConfig{}) {
		return cachedUserConfig
	}

	userConfig := getUserConfigUnvalidated()

	err := validateUserConfig(&userConfig)
	if err != nil {
		panic(err)
	}

	cachedUserConfig = userConfig
	return userConfig
}

func GetRepoConfig() RepoConfig {
	if (cachedRepoConfig != RepoConfig{}) {
		return cachedRepoConfig
	}

	userConfig := GetUserConfig()
	repoPath := userConfig.Repository.DirectoryPath

	repoPath, err := fs.AbsPath(repoPath)
	if err != nil {
		panic(err)
	}

	err = SetupRepoConfigIfNotExists(repoPath)
	if err != nil {
		slog.Error("Failed to setup repository config", "error", err)
		panic(err)
	}

	repoConfigPath := filepath.Join(repoPath, "config.toml")
	repoViper := viper.New()
	repoViper.SetConfigFile(repoConfigPath)

	err = repoViper.ReadInConfig()
	if err != nil {
		slog.Error("Error reading repository config file", "error", err, "path", repoConfigPath)
		panic(err)
	}

	var repoConfig RepoConfig
	err = repoViper.Unmarshal(&repoConfig)
	if err != nil {
		panic(err)
	}

	err = validateRepoConfig(&repoConfig, repoPath)
	if err != nil {
		panic(err)
	}

	cachedRepoPath = repoPath
	cachedRepoConfig = repoConfig

	return cachedRepoConfig
}

func GetRepoPath() string {
	if cachedRepoPath != "" {
		return cachedRepoPath
	}

	GetRepoConfig()
	return cachedRepoPath
}
