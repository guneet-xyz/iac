package utils

import (
	"iac/config"
	"iac/utils/docker"
	"iac/utils/env/secret"
	"os"
)

func SanityChecks() error {
	if exitCode, err := docker.SanityChecks(); err != nil {
		os.Exit(exitCode)
	}
	if exitCode, err := secret.SanityChecks(); err != nil {
		os.Exit(exitCode)
	}
	if exitCode, err := config.SanityChecks(); err != nil {
		os.Exit(exitCode)
	}
	return nil
}
