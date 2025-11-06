package compose

import (
	"iac/utils/os/exec"
	"log/slog"
)

func Validate(path string) (bool, error) {
	slog.Debug("Validate (Enter)", "path", path)
	cmd := exec.Command("docker", "compose", "-f", path, "config", "--quiet")
	if err := cmd.Run(); err != nil {
		slog.Debug("Validate (Exit)", "valid", false, "error", err)
		return false, err
	}
	slog.Debug("Validate (Exit)", "valid", true)
	return true, nil
}
