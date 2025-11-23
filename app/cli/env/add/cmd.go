package add

import (
	"iac/utils/env"
	"log/slog"

	"github.com/spf13/cobra"
)

var (
	envType string
	name    string
	value   string
	file    string
)

var Cmd = &cobra.Command{
	Use:     "add",
	Short:   "add a new environment secret / variable",
	Aliases: []string{"create", "new", "set"},
	RunE: func(cmd *cobra.Command, args []string) error {

		if len(name) == 0 {
			slog.Error("name option is required")
			return nil
		}
		// TODO: name validations
		// TODO: name transformations

		if len(file) > 0 {
			if len(value) > 0 {
				slog.Error("cannot use both value and file options")
				return nil
			}

			// TODO: implement reading from file
			slog.Error("not implemented")
			return nil
		} else if len(value) > 0 {
			// TODO: value validations. (is this possible?)
		} else {
			slog.Error("either value or file option must be provided")
			return nil
		}

		if envType != string(env.EnvTypeSecret) && envType != string(env.EnvTypeVariable) {
			slog.Error("type option must be either 'secret' or 'variable'")
			return nil
		}

		envData := env.Env{
			Type:  env.EnvType(envType),
			Name:  name,
			Value: value,
		}

		err := env.SetEnv(envData)
		if err != nil {
			slog.Error("failed to set environment entry", "error", err, "name", name)
			return err
		}

		slog.Info("environment entry added successfully", "name", name)
		return nil
	},
}

func init() {
	Cmd.Flags().StringVar(&envType, "type", "", "type of the environment entry (secret or variable)")
	Cmd.Flags().StringVar(&name, "name", "", "name of the secret")
	Cmd.Flags().StringVar(&value, "value", "", "value of the secret")
	Cmd.Flags().StringVar(&file, "file", "", "file containing the secret value")
}
