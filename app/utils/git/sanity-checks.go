package git

import (
	"iac/config"
	"iac/utils/errors"
	"iac/utils/exitcodes"
	"iac/utils/fs"
	"log/slog"
	"os/exec"
)

func CheckGitInstalled() error {
	slog.Debug("Checking if git is installed")

	_, err := exec.LookPath("git")
	if err != nil {
		return errors.New("git is not installed or not found in PATH", "error", err)
	}

	slog.Debug("Git is installed")
	return nil
}

func checkGitExists() (int, error) {
	slog.Debug("Checking if git is installed")

	_, err := exec.LookPath("git")
	if err != nil {
		return exitcodes.GitCommandNotFound, errors.New("git is not installed or not found in PATH", "error", err)
	}

	slog.Debug("Git is installed")
	return 0, nil
}

func checkRepositoryExists() (int, error) {
	slog.Debug("Checking if repository exists")

	cfg := config.GetUserConfig()
	repoPath := cfg.Repository.DirectoryPath

	statResult, err := fs.Stat(repoPath)
	if err != nil {
		return exitcodes.RepositoryNotFound, err
	}

	if statResult == fs.StatResultNotExist {
		return exitcodes.RepositoryNotFound, errors.New("repository directory does not exist", "path", repoPath)
	}

	if statResult != fs.StatResultDirectory {
		return exitcodes.RepositoryNotFound, errors.New("repository path is not a directory", "path", repoPath)
	}

	isRepo, err := IsGitRepository(repoPath)
	if err != nil {
		return exitcodes.RepositoryNotFound, err
	}

	if !isRepo {
		return exitcodes.RepositoryNotFound, errors.New("directory is not a git repository", "path", repoPath)
	}

	slog.Debug("Repository exists and is valid", "path", repoPath)
	return 0, nil
}

func SanityChecks() (int, error) {
	slog.Debug("Performing git sanity checks")

	if exitCode, err := checkGitExists(); err != nil {
		return exitCode, err
	}

	if exitCode, err := checkRepositoryExists(); err != nil {
		return exitCode, err
	}

	slog.Debug("Git sanity checks passed")
	return 0, nil
}
