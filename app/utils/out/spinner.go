package out

import "github.com/fatih/color"

var spinnerCharset = []string{
	"⠋",
	"⠙",
	"⠹",
	"⠸",
	"⠼",
	"⠴",
	"⠦",
	"⠧",
	"⠇",
	"⠏",
}

var spinnerIteration int

func SpinnerChar() string {
	return color.BlueString(spinnerCharset[spinnerIteration%len(spinnerCharset)])
}
