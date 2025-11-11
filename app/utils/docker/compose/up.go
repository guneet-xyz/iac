package compose

import (
	"iac/utils/os/exec"
	"log/slog"
)

func Up(composePath string, basePath string) error {
	slog.Debug("Bringing up docker compose services", "composePath", composePath, "basePath", basePath)
	cmd := exec.Command("docker", "compose", "-f", composePath, "--project-directory", basePath, "up", "-d", "--remove-orphans")
	logs, err := cmd.CombinedOutput()
	if err != nil {
		slog.Debug("Failed to bring up docker compose services", "error", err, "logs", string(logs))
		return err
	}
	slog.Debug("Successfully brought up docker compose services", "logs", string(logs))
	return nil
}
