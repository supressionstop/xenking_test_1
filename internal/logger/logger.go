package logger

import (
	"log/slog"
	"os"
	"xenking_test_1/internal/config"
)

func MustSetup(cfg config.Config) *slog.Logger {
	var handler slog.Handler
	switch cfg.App.Environment {
	default:
		handler = slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug})
		logger := slog.New(handler)
		return logger.With(slog.Group(
			"app",
			slog.String("name", cfg.App.Name),
			slog.String("env", cfg.App.Environment),
		))
	}
}
