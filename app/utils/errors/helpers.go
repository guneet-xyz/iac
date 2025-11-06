package errors

import (
	"errors"
	"fmt"
	"log/slog"
)

// New creates a new error with the given text and arguments, logs it, and returns the error.
func New(text string, args ...any) error {
	slog.Error(text, args...)

	errorMessage := text
	for i := 0; i < len(args); i += 2 {
		if i+1 < len(args) {
			errorMessage += fmt.Sprintf(" %v=%v", args[i], args[i+1])
		}
	}

	return errors.New(errorMessage)
}
