package logger

import (
	"log/slog"
	"os"

	"github.com/lmittmann/tint"
)

type LoggerOptions struct {
	Verbose bool
	Time    bool
}

func SetupLogger(opts LoggerOptions) {
	var level slog.Level
	if opts.Verbose {
		level = slog.LevelDebug
	} else {
		level = slog.LevelInfo
	}

	handlerOptions := tint.Options{
		Level: level,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			if !opts.Time && a.Key == slog.TimeKey && len(groups) == 0 {
				return slog.Attr{}
			}
			if a.Key == slog.LevelKey {
				level := a.Value.Any().(slog.Level)
				if level == slog.LevelDebug {
					return tint.Attr(244, a)
				}
			}
			return a
		},
	}

	handler := tint.NewHandler(os.Stdout, &handlerOptions)
	logger := slog.New(handler)
	slog.SetDefault(logger)
}
