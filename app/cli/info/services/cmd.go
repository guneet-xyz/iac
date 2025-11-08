package services

import (
	"fmt"
	"iac/utils/exitcodes"
	"iac/utils/service"
	"log/slog"
	"os"

	"github.com/spf13/cobra"
)

var (
	name string
)

var Cmd = &cobra.Command{
	Use:   "services",
	Short: "info about services",
	RunE:  RunE,
}

func init() {
	Cmd.Flags().StringVar(&name, "name", "", "Specify a service to get info about")
}

func RunE(cmd *cobra.Command, args []string) error {
	services, err := service.GetServiceNames()
	if err != nil {
		return err
	}

	var filteredServices []string

	if name != "" {
		for _, svc := range services {
			if svc == name {
				filteredServices = append(filteredServices, svc)
				break
			}
		}
		if len(filteredServices) == 0 {
			slog.Info("Service not found", "name", name)
			os.Exit(exitcodes.ServiceNotFound)
			return nil
		}
	} else {
		filteredServices = services
		if len(filteredServices) == 0 {
			slog.Info("No services found")
			os.Exit(exitcodes.NoServicesFound)
			return nil
		}
	}

	for _, svc := range filteredServices {
		PrintServiceInfo(svc)
	}

	return nil
}

func PrintServiceInfo(svcName string) error {
	info, err := service.GetServiceInfo(svcName)
	if err != nil {
		return err
	}

	fmt.Printf("Service: %s\n", info.Name)
	for _, container := range info.Containers {
		fmt.Printf("  Container Name: %s\n", container.ConfigInfo.ContainerName)
		if container.ConfigInfoFound {
			fmt.Printf("    Config Info Found: Yes\n")
		} else {
			fmt.Printf("    Config Info Found: No\n")
		}
		if container.RunningInfoFound {
			fmt.Printf("    Running Info Found: Yes\n")
			fmt.Printf("    Container ID: %s\n", container.RunningInfo.ContainerId)
		} else {
			fmt.Printf("    Running Info Found: No\n")
		}
	}
	return nil
}
