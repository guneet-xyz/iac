package out

import (
	"fmt"
	"iac/utils/ansi"
)

var linesPreviouslyPrinted int

// clears the previously printed lines and prints the new lines.
func RepaintLines(lines []string) {
	for i := 0; i < linesPreviouslyPrinted; i++ {
		fmt.Print(ansi.LineUp + ansi.LineClear)
	}

	linesPreviouslyPrinted = len(lines)

	for _, line := range lines {
		fmt.Println(line)
	}

	spinnerIteration++
}
