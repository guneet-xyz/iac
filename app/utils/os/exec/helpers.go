package exec

import (
	"fmt"
	"log/slog"
	exec "os/exec"
)

func Command(name string, arg ...string) *exec.Cmd {
	cmd := fmt.Sprintf("%s %v", name, arg)
	slog.Debug("Running command", "cmd", cmd)
	return exec.Command(name, arg...)
}
