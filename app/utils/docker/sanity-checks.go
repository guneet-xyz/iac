package docker

import (
	"iac/utils/errors"
	"iac/utils/exitcodes"
	"log/slog"
	"os/exec"
)

func exists() (int, error) {
	slog.Debug("Checking if Docker is installed")
	cmd := exec.Command("docker", "version")
	err := cmd.Run()
	if err != nil {
		return exitcodes.DockerCommandNotFound, errors.New("Docker is not installed", "error", err)
	}
	exitCode := cmd.ProcessState.ExitCode()
	if exitCode != 0 {
		return exitcodes.DockerCommandNotFound, errors.New("Docker is not installed", "exit_code", exitCode)
	}
	slog.Debug("Docker is installed")
	return 0, nil
}

func havePermission() (int, error) {
	slog.Debug("Checking if user has permission to run Docker commands")
	cmd := exec.Command("docker", "ps")
	err := cmd.Run()
	if err != nil {
		return exitcodes.DockerInsufficientPermissions, errors.New("Insufficient permissions to run Docker commands", "error", err)
	}
	exitCode := cmd.ProcessState.ExitCode()
	if exitCode != 0 {
		return exitcodes.DockerInsufficientPermissions, errors.New("Insufficient permissions to run Docker commands", "exit_code", exitCode)
	}
	slog.Debug("User has permission to run Docker commands")
	return 0, nil
}

func SanityChecks() (int, error) {
	slog.Debug("Performing Docker sanity checks")

	if exitCode, err := exists(); err != nil {
		return exitCode, err
	}

	if exitCode, err := havePermission(); err != nil {
		return exitCode, err
	}

	slog.Debug("Docker sanity checks passed")
	return 0, nil
}
