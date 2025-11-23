package get

import (
	"fmt"
	"iac/utils/env"
	"log/slog"

	"github.com/spf13/cobra"
)

var (
	name string
)

var Cmd = &cobra.Command{
	Use:   "get",
	Short: "get an environment entry",
	RunE: func(cmd *cobra.Command, args []string) error {
		envData, err := env.GetEnv(name)
		if err != nil {
			slog.Error("Failed to get environment entry", "error", err, "name", name)
			return err
		}
		fmt.Println(envData.Value)
		return nil
	},
}

func init() {
	Cmd.Flags().StringVarP(&name, "name", "n", "", "name of the secret")
}
