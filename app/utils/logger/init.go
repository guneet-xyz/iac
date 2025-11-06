// This package exists to that this logger setup happens before any other package code runs
package logger

import (
	"log/slog"
	"os"

	"github.com/lmittmann/tint"
)

func init() {
	handlerOptions := tint.Options{
		Level: slog.LevelInfo,
	}
	handler := tint.NewHandler(os.Stdout, &handlerOptions)
	logger := slog.New(handler)
	slog.SetDefault(logger)
}
