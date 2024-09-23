package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"xenking_test_1/internal/config"
	internalLogger "xenking_test_1/internal/logger"
	"xenking_test_1/internal/provider"
	"xenking_test_1/internal/server"
	"xenking_test_1/internal/storage"
	"xenking_test_1/internal/usecase"
	"xenking_test_1/internal/usecase/repository"
	"xenking_test_1/internal/worker"
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

	// storage
	// TODO: separate sports to tables
	pg, err := storage.NewPostgres(ctx, cfg.DB.URL, logger)
	if err != nil {
		return err
	}
	defer pg.ClosePool()

	// migrations
	if err := pg.Up(cfg.DB.URL); err != nil {
		return err
	}

	// provider
	kiddy, err := provider.NewKiddy(&http.Client{Timeout: cfg.Provider.HttpTimeout}, cfg.Provider.BaseUrl)
	if err != nil {
		return err
	}

	// domain
	lineRepository := repository.NewLine(pg)
	getLine := usecase.NewGetLineUseCase(kiddy)
	saveLine := usecase.NewSaveLineUseCase(lineRepository)
	fetchLine := usecase.NewFetchLineUseCase(getLine, saveLine)

	// workers
	worker.NewPool(*cfg, logger, fetchLine).StartWorkers(ctx)

	// http server
	server.NewHTTP(cfg, logger, workerPool).Start(ctx)

	logger.Info("running")
	select {
	case <-ctx.Done():
		stop()
		logger.Info("shutting down")
	}

	return nil
}
