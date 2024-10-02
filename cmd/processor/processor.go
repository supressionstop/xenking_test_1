package main

import (
	"context"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	grpcController "github.com/supressionstop/xenking_test_1/internal/application/controller/grpc"
	httpController "github.com/supressionstop/xenking_test_1/internal/application/controller/http"
	"github.com/supressionstop/xenking_test_1/internal/application/task"
	"github.com/supressionstop/xenking_test_1/internal/infrastructure/config"
	internalLogger "github.com/supressionstop/xenking_test_1/internal/infrastructure/logger"
	"github.com/supressionstop/xenking_test_1/internal/infrastructure/provider"
	"github.com/supressionstop/xenking_test_1/internal/infrastructure/server"
	"github.com/supressionstop/xenking_test_1/internal/infrastructure/storage"
	"github.com/supressionstop/xenking_test_1/internal/infrastructure/subscription"
	"github.com/supressionstop/xenking_test_1/internal/infrastructure/worker"
	"github.com/supressionstop/xenking_test_1/internal/usecase"
	"github.com/supressionstop/xenking_test_1/internal/usecase/repository"
)

func main() {
	ctx := context.Background()

	ctx, stop := signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	cfg, err := config.Setup(os.Getenv("APP_ENV"))
	if err != nil {
		log.Fatalf("failed to setup config: %v", err)
	}

	logger := internalLogger.MustSetup(*cfg)
	logger.Info("starting...")

	// storage
	storageInstance, err := storage.NewPostgres(ctx, cfg.DB.URL, logger)
	if err != nil {
		logger.Error("failed to setup storage", slog.Any("error", err))
		os.Exit(1)
	}
	defer storageInstance.ClosePool()

	// migrations
	if err := storageInstance.Up(cfg.DB.URL); err != nil {
		logger.Error("failed to apply migrations", slog.Any("error", err))
		os.Exit(1)
	}

	// provider
	kiddy, err := provider.NewKiddy(cfg.Provider.BaseURL, cfg.Provider.HTTPTimeout)
	if err != nil {
		logger.Error("failed to setup provider", slog.Any("error", err))
		os.Exit(1)
	}

	// domain
	lineRepository := repository.NewLine(storageInstance)
	getLine := usecase.NewGetLineUseCase(kiddy)
	saveLine := usecase.NewSaveLineUseCase(lineRepository)
	fetchLine := usecase.NewFetchLineUseCase(getLine, saveLine)
	getRecentLines := usecase.NewGetRecentLinesUseCase(lineRepository)
	calculateDiff := usecase.NewCalculateDiffUseCase()

	// worker pool
	taskGenerator := task.NewGenerator(cfg.Workers, fetchLine)
	taskGenerator.Start(ctx)
	workerPool := worker.NewPool(len(cfg.Workers), taskGenerator.Tasks(), logger)
	go workerPool.Start(ctx)

	// domain
	isLineSynced := usecase.NewIsLineSynced(workerPool)

	// http server
	readyController := httpController.NewReadyController(isLineSynced)
	httpServer := server.NewHttp(cfg, logger, readyController)
	go httpServer.Start(ctx)

	// grpc server
	subscriptionManager := subscription.NewSubscriptionManager(getRecentLines, calculateDiff, logger)
	subscriptionController := grpcController.NewSubscriptionController(subscriptionManager, logger)
	grpcServer := server.NewGrpc(
		":"+cfg.GRPCServer.Port,
		logger,
		subscriptionController,
	)
	go grpcServer.DeferredStart(ctx, cfg.MaxWorkerInterval(), workerPool.AllWorkersSynced)

	logger.Info("running")

	<-ctx.Done()
	stop()
	<-httpServer.Finish
	<-grpcServer.Finish
	logger.Info("app finished.")
}
