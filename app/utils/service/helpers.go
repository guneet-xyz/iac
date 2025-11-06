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

func GetServiceNames() ([]string, error) {
	slog.Debug("Getting service names from services directory")

	conf := config.GetConfig()

	services := []string{}

	err := filepath.WalkDir(conf.ServicesDirPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		var composeFound bool = false
		if d.IsDir() && path != conf.ServicesDirPath {
			for _, fname := range possibleComposeFilenames {
				stat, err :=
					fs.Stat(filepath.Join(path, fname))
				if err != nil {
					return err
				}
				if stat != fs.StatResultFile {
					continue
				}
				absPath, err := filepath.Abs(filepath.Join(path, fname))
				if err != nil {
					return err
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
				composeFound = true
				break
			}

			if composeFound {
				services = append(services, d.Name())
			}

			return filepath.SkipDir
		}

		return nil
	})

	if err != nil {
		slog.Debug("Error reading services directory: %v", err)
		return nil, err
	}

	return services, nil
}
