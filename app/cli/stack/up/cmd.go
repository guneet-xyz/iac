package up

import (
	"iac/cli/stack/info"
	"iac/utils/docker/compose"
	"iac/utils/out"
	"iac/utils/stack"
	"log/slog"
	"path/filepath"
	"sync"
	"time"

	"github.com/spf13/cobra"
)

var (
	stackName string
)

var Cmd = &cobra.Command{
	Use:   "up",
	Short: "bring up stacks",
	RunE: func(_ *cobra.Command, _ []string) error {
		if stackName == "" {
			slog.Error("stack name is required")
			return nil
		}

		runSpinner := true
		wg := sync.WaitGroup{}

		wg.Go(func() {
			for runSpinner {
				lines := []string{
					out.SpinnerChar(),
				}
				out.RepaintLines(lines)
				time.Sleep(100 * time.Millisecond)
			}
		})

		composePath, err := stack.GetComposePathFromStackName(stackName)
		if err != nil {
			return err
		}
		composeDir := filepath.Dir(composePath)
		err = compose.Up(composePath, composeDir)
		if err != nil {
			return err
		}

		runSpinner = false
		wg.Wait()

		err = info.ShowInfoAboutStack(stackName)
		return err
	},
}

func init() {
	Cmd.Flags().StringVarP(&stackName, "stack", "s", "", "Name of the stack to bring up")
}
