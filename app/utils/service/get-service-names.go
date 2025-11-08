package service

import (
	"iac/config"
	"iac/utils/fs"
	"log/slog"
	"path/filepath"
)

func GetServiceNames() ([]string, error) {
	slog.Debug("Getting service names from services directory")

	conf := config.GetConfig()

	services := []string{}

	err := filepath.WalkDir(conf.ServicesDirPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if path == conf.ServicesDirPath {
			return nil
		}

		composePath, err := GetComposePathFromServiceName(d.Name())
		if composePath != "" && err == nil {
			services = append(services, d.Name())
		}

		return filepath.SkipDir
	})

	if err != nil {
		slog.Debug("Error reading services directory", "error", err)
		return nil, err
	}

	return services, nil
}
