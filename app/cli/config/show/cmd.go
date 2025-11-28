package show

import (
	"encoding/json"
	"fmt"
	"log/slog"

	"iac/config"
	"iac/utils/errors"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	jsonOutput bool
)

var Cmd = &cobra.Command{
	Use:   "show",
	Short: "show configuration settings",
	RunE: func(cmd *cobra.Command, args []string) error {
		config := config.GetConfig()
		configFile := viper.GetViper().ConfigFileUsed()

		if jsonOutput {
			slog.Debug("Json Output is enabled")

			jsonData, err := json.MarshalIndent(config, "", "  ")
			if err != nil {
				return errors.New("Failed to generate JSON output", "error", err)
			}
			fmt.Println(string(jsonData))

		} else {
			slog.Debug("Standard Output is enabled")

			fmt.Printf("Configuration file is at %s\n", configFile)

			fmt.Printf("Repository Directory: %s\n", config.Repository.DirectoryPath)
			fmt.Printf("Repository Origin: %s\n", config.Repository.OriginURL)
			fmt.Printf("Environment Directory: %s\n", config.Environment.DirectoryPath)
			fmt.Printf("Stacks Directory: %s\n", config.Stacks.DirectoryPath)
			fmt.Printf("Backups Directory: %s\n", config.Backups.DirectoryPath)
			fmt.Printf("Master Key Path: %s\n", config.MasterKey.MasterKeyPath)
			fmt.Printf("KCV Path: %s\n", config.MasterKey.KcvPath)
		}

		return nil
	},
}

func init() {
	Cmd.Flags().BoolVarP(&jsonOutput, "json", "j", false, "Output in JSON format")
}
