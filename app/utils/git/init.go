package git

import (
	_ "embed"
	"iac/config"
	"log/slog"
	"path/filepath"

	"iac/utils/errors"
	"iac/utils/fs"
	"iac/utils/os/exec"
)

//go:embed templates/README.md
var readmeTemplate string

func InitRepository(path string) error {
	slog.Debug("Initializing repository", "path", path)

	statResult, err := fs.Stat(path)
	if err != nil {
		return err
	}
	if statResult == fs.StatResultNotExist {
		slog.Debug("Creating directory", "path", path)
		expandedPath, err := fs.MkdirIfNotExists(path)
		if err != nil {
			return errors.New("failed to create directory", "path", path, "error", err)
		}
		path = expandedPath
	} else if statResult == fs.StatResultFile {
		return errors.New("path exists but is a file, not a directory", "path", path)
	}

	isRepo, err := IsGitRepository(path)
	if err != nil {
		return err
	}
	if isRepo {
		return errors.New("directory is already a git repository", "path", path)
	}

	cmd := exec.Command("git", "-C", path, "init")
	output, err := cmd.CombinedOutput()
	if err != nil {
		slog.Debug("Failed to initialize repository", "error", err, "output", string(output))
		return errors.New("failed to initialize git repository", "path", path, "error", err, "output", string(output))
	}

	slog.Debug("Git repository initialized", "output", string(output))

	defaultDirs := []string{"stacks", "env", "backups"}
	for _, dir := range defaultDirs {
		dirPath := filepath.Join(path, dir)
		_, err := fs.MkdirIfNotExists(dirPath)
		if err != nil {
			slog.Warn("Failed to create directory", "path", dirPath, "error", err)
			continue
		}
		slog.Info("Created directory", "path", dirPath)

		gitkeepPath := filepath.Join(dirPath, ".gitkeep")
		err = fs.WriteFileFromString(gitkeepPath, "")
		if err != nil {
			slog.Warn("Failed to create .gitkeep file", "path", gitkeepPath, "error", err)
		}
	}

	readmePath := filepath.Join(path, "README.md")
	err = fs.WriteFileFromString(readmePath, readmeTemplate)
	if err != nil {
		slog.Warn("Failed to create README.md", "path", readmePath, "error", err)
	}

	err = config.SetupRepoConfigIfNotExists(path)
	if err != nil {
		slog.Warn("Failed to create repository config", "error", err)
	}

	slog.Info("Repository initialized with default structure", "path", path)
	return nil
}
