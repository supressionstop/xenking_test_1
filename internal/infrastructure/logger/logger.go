package logger

import (
	"log"
	"log/slog"
	"os"

	"github.com/supressionstop/xenking_test_1/internal/infrastructure/config"
)

func MustSetup(cfg config.Config) *slog.Logger {
	var handler slog.Handler

	level, err := parseLevel(cfg.Log.Level)
	if err != nil {
		log.Fatal(err)
	}

	switch cfg.App.Environment {
	default:
		handler = slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: level})
		logger := slog.New(handler)
		return logger.With(slog.Group(
			"app",
			slog.String("name", cfg.App.Name),
			slog.String("env", cfg.App.Environment),
		))
	}
}

func parseLevel(s string) (slog.Level, error) {
	var level slog.Level
	err := level.UnmarshalText([]byte(s))

	return level, err
}
