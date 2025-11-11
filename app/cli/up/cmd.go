package up

import (
	"fmt"
	"iac/cli/info/services"
	"iac/utils/docker/compose"
	"iac/utils/service"
	"log/slog"
	"path/filepath"

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
		composePath, err := service.GetComposePathFromServiceName(serviceName)
		if err != nil {
			return err
		}
		composeDir := filepath.Dir(composePath)
		err = compose.Up(composePath, composeDir)
		if err != nil {
			return err
		}
		err = services.ShowInfoAboutService(serviceName)
		return err
	},
}

func init() {
	Cmd.Flags().StringVarP(&serviceName, "service", "s", "", "Name of the service to bring up")
}
