package main

import (
	"iac/cli"
	"iac/utils/exitcodes"
	_ "iac/utils/logger"
	"os"
)

func main() {
	err := cli.Cmd.Execute()
	if err != nil {
		os.Exit(exitcodes.UnknownError)
	}
}
