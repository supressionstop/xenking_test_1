package main

import (
	"log"
	"log/slog"
	"os"
	"xenking_test_1/internal/config"
	internalLogger "xenking_test_1/internal/logger"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	cfg, err := config.MustSetup(os.Getenv("APP_ENVIRONMENT"))
	if err != nil {
		return err
	}

	logger := internalLogger.Must(cfg.Environment)
	logger.Info("starting...", slog.String("name", cfg.Name), slog.String("env", cfg.Environment))

	// graceful shutdown
	// line provider client
	// storage
	// workers

	return nil
}
