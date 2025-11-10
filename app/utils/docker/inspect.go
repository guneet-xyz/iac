package docker

import (
	"iac/utils/json"
	"iac/utils/os/exec"
	"log/slog"
	"time"
)

type InspectResultState struct {
	Running   bool
	StartedAt time.Time
}
type InspectResult struct {
	Id           string
	Name         string
	RestartCount int
	State        InspectResultState
}

func Inspect(containerSelector string) ([]InspectResult, error) {
	slog.Debug("Inspecting container", "selector", containerSelector)

	cmd := exec.Command("docker", "inspect", containerSelector)
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}
	slog.Debug("Docker inspect output", "output", string(output))

	var unmarshaled []InspectResult
	err = json.Unmarshal(output, &unmarshaled)
	if err != nil {
		return nil, err
	}
	slog.Debug("Unmarshaled inspect result", "result", unmarshaled)

	return unmarshaled, nil
}
