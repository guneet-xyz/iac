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
		userConfig := config.GetUserConfig()
		repoConfig := config.GetRepoConfig()
		configFile := viper.GetViper().ConfigFileUsed()

		if jsonOutput {
			slog.Debug("Json Output is enabled")

			output := map[string]interface{}{
				"user":       userConfig,
				"repository": repoConfig,
			}
			jsonData, err := json.MarshalIndent(output, "", "  ")
			if err != nil {
				return errors.New("Failed to generate JSON output", "error", err)
			}
			fmt.Println(string(jsonData))

		} else {
			slog.Debug("Standard Output is enabled")

			fmt.Printf("Configuration file: %s\n\n", configFile)

			fmt.Println("Repository:")
			fmt.Printf("  Directory: %s\n\n", userConfig.Repository.DirectoryPath)

			fmt.Println("Key:")
			fmt.Printf("  Directory: %s\n\n", repoConfig.Key.DirectoryPath)

			fmt.Println("Directories:")
			fmt.Printf("  Environment: %s\n", repoConfig.Environment.DirectoryPath)
			fmt.Printf("  Stacks:      %s\n", repoConfig.Stacks.DirectoryPath)
			fmt.Printf("  Backups:     %s\n", repoConfig.Backups.DirectoryPath)
		}

		return nil
	},
}

func init() {
	Cmd.Flags().BoolVarP(&jsonOutput, "json", "j", false, "Output in JSON format")
}
