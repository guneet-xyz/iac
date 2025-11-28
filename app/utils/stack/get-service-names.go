package stack

import (
	"iac/config"
	"iac/utils/fs"
	"log/slog"
	"path/filepath"
)

func GetStackNames() ([]string, error) {
	slog.Debug("Getting stack names from stacks directory")

	conf := config.GetConfig()

	stacks := []string{}

	err := filepath.WalkDir(conf.Stacks.DirectoryPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if path == conf.Stacks.DirectoryPath {
			return nil
		}

		composePath, err := GetComposePathFromStackName(d.Name())
		if composePath != "" && err == nil {
			stacks = append(stacks, d.Name())
		}

		return filepath.SkipDir
	})

	if err != nil {
		slog.Debug("Error reading stacks directory", "error", err)
		return nil, err
	}

	return stacks, nil
}
