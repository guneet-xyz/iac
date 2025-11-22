package logger

import (
	"log/slog"
	"os"

	"github.com/lmittmann/tint"
)

func SetupLogger(verbose bool) {
	var level slog.Level
	if verbose {
		level = slog.LevelDebug
	} else {
		level = slog.LevelInfo
	}
	handlerOptions := tint.Options{
		Level: level,
	}
	handler := tint.NewHandler(os.Stdout, &handlerOptions)
	logger := slog.New(handler)
	slog.SetDefault(logger)
}
