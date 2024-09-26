package main

import (
	"context"
	"fmt"
	"github.com/supressionstop/xenking_test_1/internal/provider"
	"github.com/supressionstop/xenking_test_1/internal/server"
	"github.com/supressionstop/xenking_test_1/internal/usecase"
	"github.com/supressionstop/xenking_test_1/internal/usecase/repository"
	"github.com/supressionstop/xenking_test_1/internal/worker"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/supressionstop/xenking_test_1/internal/config"
	internalLogger "github.com/supressionstop/xenking_test_1/internal/logger"
	"github.com/supressionstop/xenking_test_1/internal/storage"
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
		return fmt.Errorf("failed to setup config: %w", err)
	}

	logger := internalLogger.MustSetup(*cfg)
	logger.Info("starting...")

	// storage
	storageInstance, err := storage.NewPostgres(ctx, cfg.DB.URL, logger)
	if err != nil {
		return fmt.Errorf("failed to setup storage: %w", err)
	}
	defer storageInstance.ClosePool()

	// migrations
	if err := storageInstance.Up(cfg.DB.URL); err != nil {
		return fmt.Errorf("failed to apply migrations: %w", err)
	}

	// provider
	httpClient := &http.Client{
		Transport:     nil,
		CheckRedirect: nil,
		Jar:           nil,
		Timeout:       cfg.Provider.HTTPTimeout,
	}

	kiddy, err := provider.NewKiddy(httpClient, cfg.Provider.BaseURL)
	if err != nil {
		return fmt.Errorf("failed to setup provider: %w", err)
	}

	// domain
	lineRepository := repository.NewLine(storageInstance)
	getLine := usecase.NewGetLineUseCase(kiddy)
	saveLine := usecase.NewSaveLineUseCase(lineRepository)
	fetchLine := usecase.NewFetchLineUseCase(getLine, saveLine)
	getRecentLines := usecase.NewGetRecentLinesUseCase(lineRepository)
	calculateDiff := usecase.NewCalculateDiffUseCase()

	// workers
	workerPool := worker.NewPool(*cfg, logger, fetchLine)
	workerPool.StartWorkers(ctx)

	// http server
	httpServer := server.NewHttp(cfg, logger, workerPool)
	httpServer.Start()

	// grpc server
	subscriptionManager := server.NewSubscriptionManager(getRecentLines, calculateDiff, logger)
	grpcServer := server.NewGrpc(
		":"+cfg.GRPCServer.Port,
		logger,
		subscriptionManager,
	)

	err = grpcServer.DeferredStart(workerPool.FirstSyncChan)
	if err != nil {
		return fmt.Errorf("failed to start gRPC server: %w", err)
	}

	logger.Info("running")
	select {
	case err := <-httpServer.ErrChan:
		if err != nil {
			logger.Error("http server error", slog.Any("err", err))
		}

		ctx.Done()
	case err := <-grpcServer.ErrChan:
		if err != nil {
			logger.Error("grpc server error", slog.Any("err", err))
		}

		ctx.Done()
	case <-ctx.Done():
		grpcServer.GracefulStop()

		err := httpServer.Shutdown(ctx)
		if err != nil {
			logger.Error("failed to shutdown http server", slog.Any("err", err))
		}

		stop()
	}
	logger.Info("app finished.")

	return nil
}
