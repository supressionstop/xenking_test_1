package logger

import (
	"log/slog"
	"os"
)

func Must(env string) *slog.Logger {
	var handler slog.Handler
	switch env {
	default:
		handler = slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug})
	}

	return slog.New(handler)
}
