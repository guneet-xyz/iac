package compose

import (
	"bytes"
	"encoding/json"
	"iac/utils/os/exec"
	"log/slog"
)

type StatsContainer struct {
	Id   string `json:"Container"`
	Name string `json:"Name"`
}

type Stats struct {
	Containers []StatsContainer
}

func GetStats(composePath string) (Stats, error) {
	cmd := exec.Command("docker", "compose", "-f", composePath, "stats", "--no-stream", "--format", "json")
	output, err := cmd.Output()
	if err != nil {
		slog.Error("Failed to get docker compose stats", "error", err)
		return Stats{}, err
	}

	stats := Stats{
		Containers: []StatsContainer{},
	}

	byteReader := bytes.NewReader(output)
	decoder := json.NewDecoder(byteReader)

	for decoder.More() {
		var container StatsContainer
		err = decoder.Decode(&container)
		if err != nil {
			slog.Error("Failed to decode docker compose stats JSON", "error", err)
			return Stats{}, err
		}
		stats.Containers = append(stats.Containers, container)
	}

	return stats, nil
}
