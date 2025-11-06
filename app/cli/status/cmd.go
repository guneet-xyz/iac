package status

import (
	"fmt"
	"iac/utils/service"

	"github.com/spf13/cobra"
)

var Cmd = &cobra.Command{
	Use:   "status",
	Short: "status",
	Run: func(cmd *cobra.Command, args []string) {
		services, err := service.GetServiceNames()
		if err != nil {
			fmt.Println("Error retrieving services:", err)
			return
		}

		fmt.Printf("Found %d services:\n", len(services))
		for _, svc := range services {
			fmt.Printf("- %s\n", svc)
		}
	},
}

// check config status
// check how many services are found
// check how many are deployed
