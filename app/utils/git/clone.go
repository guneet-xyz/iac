package git

import (
	"log/slog"

	"iac/utils/errors"
	"iac/utils/fs"
	"iac/utils/os/exec"
)

func CloneRepository(url string, path string) error {
	slog.Debug("Cloning repository", "url", url, "path", path)

	statResult, err := fs.Stat(path)
	if err != nil {
		return err
	}
	if statResult == fs.StatResultDirectory {
		return errors.New("destination directory already exists", "path", path)
	} else if statResult == fs.StatResultFile {
		return errors.New("destination path exists but is a file, not a directory", "path", path)
	}

	cmd := exec.Command("git", "clone", url, path)
	output, err := cmd.CombinedOutput()
	if err != nil {
		slog.Debug("Failed to clone repository", "error", err, "output", string(output))
		return errors.New("failed to clone repository", "url", url, "path", path, "error", err, "output", string(output))
	}

	slog.Info("Repository cloned successfully", "path", path)
	slog.Debug("Clone output", "output", string(output))
	return nil
}
