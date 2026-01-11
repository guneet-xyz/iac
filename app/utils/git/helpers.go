package git

import (
	"log/slog"
	"path/filepath"
	"strings"

	"iac/utils/errors"
	"iac/utils/fs"
	"iac/utils/os/exec"
)

func IsGitRepository(path string) (bool, error) {
	slog.Debug("Checking if path is a git repository", "path", path)

	statResult, err := fs.Stat(path)
	if err != nil {
		return false, err
	}
	if statResult == fs.StatResultNotExist {
		return false, nil
	}
	if statResult != fs.StatResultDirectory {
		return false, errors.New("path is not a directory", "path", path)
	}

	gitDir := filepath.Join(path, ".git")
	gitStatResult, err := fs.Stat(gitDir)
	if err != nil {
		return false, err
	}

	if gitStatResult == fs.StatResultDirectory {
		slog.Debug("Found .git directory", "path", gitDir)
		return true, nil
	}

	return false, nil
}

func GetOriginURL(path string) (string, error) {
	slog.Debug("Getting origin URL", "path", path)

	cmd := exec.Command("git", "-C", path, "config", "--get", "remote.origin.url")
	output, err := cmd.Output()
	if err != nil {
		slog.Warn("No origin URL configured for repository", "path", path)
		return "", nil
	}

	originURL := strings.TrimSpace(string(output))
	slog.Debug("Got origin URL", "url", originURL)
	return originURL, nil
}

func ValidateRepository(path string, expectedOriginURL string) error {
	slog.Debug("Validating repository", "path", path, "expectedOrigin", expectedOriginURL)

	isRepo, err := IsGitRepository(path)
	if err != nil {
		return err
	}
	if !isRepo {
		return errors.New("directory is not a git repository", "path", path)
	}

	if expectedOriginURL == "" {
		slog.Debug("No expected origin URL provided, skipping origin validation")
		return nil
	}

	originURL, err := GetOriginURL(path)
	if err != nil {
		return err
	}

	if originURL != expectedOriginURL {
		return errors.New("repository origin URL does not match expected URL", "path", path, "expected", expectedOriginURL, "actual", originURL)
	}

	slog.Debug("Repository validation successful")
	return nil
}
