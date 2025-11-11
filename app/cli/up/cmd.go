package up

import (
	"iac/cli/info/services"
	"iac/utils/docker/compose"
	"iac/utils/out"
	"iac/utils/service"
	"log/slog"
	"path/filepath"
	"sync"
	"time"

	"github.com/spf13/cobra"
)

var (
	serviceName string
)

var Cmd = &cobra.Command{
	Use:   "up",
	Short: "bring up services",
	RunE: func(_ *cobra.Command, _ []string) error {
		if serviceName == "" {
			slog.Error("Service name is required")
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

		composePath, err := service.GetComposePathFromServiceName(serviceName)
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

		err = services.ShowInfoAboutService(serviceName)
		return err
	},
}

func init() {
	Cmd.Flags().StringVarP(&serviceName, "service", "s", "", "Name of the service to bring up")
}
