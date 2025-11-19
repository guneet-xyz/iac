package utils

import (
	"iac/utils/docker"
	"iac/utils/secret"
	"os"
)

func SanityChecks() error {
	if exitCode, err := docker.SanityChecks(); err != nil {
		os.Exit(exitCode)
	}
	if exitCode, err := secret.SanityChecks(); err != nil {
		os.Exit(exitCode)
	}
	return nil
}
