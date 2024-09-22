package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"xenking_test_1/internal/config"
	internalLogger "xenking_test_1/internal/logger"
	"xenking_test_1/internal/server"
)

func main() {
	ctx := context.Background()
	if err := run(ctx); err != nil {
		log.Fatal(err)
	}
}

func run(ctx context.Context) error {
	ctx, stop := signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	cfg, err := config.MustSetup(os.Getenv("APP_ENVIRONMENT"))
	if err != nil {
		return err
	}

	logger := internalLogger.MustSetup(*cfg)
	logger.Info("starting...")

	// line provider client
	// storage
	// workers

	server.NewHTTP(cfg, logger).Start(ctx)

	logger.Info("running")
	select {
	case <-ctx.Done():
		stop()
		logger.Info("shutting down")
	}

	return nil
}
