package stack

import (
	"iac/config"
	"iac/utils/docker/compose"
	"iac/utils/fs"
	"log/slog"
	"path/filepath"
)

var possibleComposeFilenames = []string{
	"docker-compose.yml",
	"docker-compose.yaml",
	"compose.yml",
	"compose.yaml",
}

func GetStackDirFromStackName(stackName string) (string, error) {
	cfg, err := config.GetConfig()
	if err != nil {
		slog.Error("Failed to get config", "error", err)
		return "", err
	}
	return filepath.Join(cfg.Stacks.DirectoryPath, stackName), nil
}

func GetComposePathFromStackName(stackName string) (string, error) {
	slog.Debug("Searching for compose file", "stackName", stackName)
	stackDir, err := GetStackDirFromStackName(stackName)
	if err != nil {
		return "", err
	}

	for _, fname := range possibleComposeFilenames {
		stat, err :=
			fs.Stat(filepath.Join(stackDir, fname))
		if err != nil {
			return "", err
		}
		if stat != fs.StatResultFile {
			continue
		}
		absPath, err := filepath.Abs(filepath.Join(stackDir, fname))
		if err != nil {
			return "", err
		}
		valid, err := compose.Validate(absPath)
		if err != nil {
			slog.Debug("Compose file validation error", "file", absPath, "error", err)
			continue
		}
		if !valid {
			slog.Debug("Compose file is not valid", "file", absPath)
			continue
		}
		return absPath, nil
	}

	slog.Debug("No valid compose file found", "stackName", stackName)
	return "", nil
}
