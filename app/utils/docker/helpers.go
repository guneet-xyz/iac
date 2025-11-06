package docker

import (
	"iac/utils/exitcodes"
	"log/slog"
	"os"
	"os/exec"
)

func exists() bool {
	cmd := exec.Command("docker", "version")
	err := cmd.Run()
	return err == nil && cmd.ProcessState.ExitCode() == 0
}

func havePermission() bool {
	cmd := exec.Command("docker", "ps")
	err := cmd.Run()
	return err == nil && cmd.ProcessState.ExitCode() == 0
}

func SanityChecks() {
	slog.Debug("Performing Docker sanity checks")

	if !exists() {
		slog.Error("Docker is not installed or not found in PATH")
		os.Exit(exitcodes.DockerCommandNotFound)
	}

	if !havePermission() {
		os.Exit(exitcodes.DockerInsufficientPermissions)
		slog.Error("Insufficient permissions to run Docker commands")
	}
}

func init() {
	SanityChecks()
}
