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
		config, err := config.GetConfig()
		if err != nil {
			return errors.New("Failed to retrieve configuration", "error", err)
		}
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

			fmt.Printf("Configuration file: %s\n\n", configFile)

			fmt.Println("Repository:")
			fmt.Printf("  Directory: %s\n", config.Repository.DirectoryPath)
			fmt.Printf("  Origin:    %s\n\n", config.Repository.OriginURL)

			fmt.Println("Master Key:")
			fmt.Printf("  Key Path: %s\n", config.MasterKey.MasterKeyPath)
			fmt.Printf("  KCV Path: %s\n\n", config.MasterKey.KcvPath)

			fmt.Println("Directories:")
			fmt.Printf("  Environment: %s\n", config.Environment.DirectoryPath)
			fmt.Printf("  Stacks:      %s\n", config.Stacks.DirectoryPath)
			fmt.Printf("  Backups:     %s\n", config.Backups.DirectoryPath)
		}

		return nil
	},
}

func init() {
	Cmd.Flags().BoolVarP(&jsonOutput, "json", "j", false, "Output in JSON format")
}
