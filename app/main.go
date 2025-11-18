package main

import (
	"iac/cli"
	"iac/utils/exitcodes"
	"iac/utils/logger"
	"os"
)

func main() {
	logger.SetupLogger()
	err := cli.Cmd.Execute()
	if err != nil {
		os.Exit(exitcodes.UnknownError)
	}
}
