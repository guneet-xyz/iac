package service

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

func GetServiceDirFromServiceName(svcName string) string {
	return filepath.Join(config.GetConfig().ServicesDirPath, svcName)
}

func GetComposePathFromServiceName(svcName string) (string, error) {
	slog.Debug("Searching for compose file", "service", svcName)
	svcDir := GetServiceDirFromServiceName(svcName)

	for _, fname := range possibleComposeFilenames {
		stat, err :=
			fs.Stat(filepath.Join(svcDir, fname))
		if err != nil {
			return "", err
		}
		if stat != fs.StatResultFile {
			continue
		}
		absPath, err := filepath.Abs(filepath.Join(svcDir, fname))
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

	slog.Debug("No valid compose file found", "service", svcName)
	return "", nil
}
