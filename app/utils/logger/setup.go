package logger

import (
	"log/slog"
	"os"

	"github.com/lmittmann/tint"
)

func SetupLogger() {
	handlerOptions := tint.Options{
		Level: slog.LevelInfo,
	}
	handler := tint.NewHandler(os.Stdout, &handlerOptions)
	logger := slog.New(handler)
	slog.SetDefault(logger)
}
