package server

import (
	"context"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	slogchi "github.com/samber/slog-chi"
	"github.com/supressionstop/xenking_test_1/internal/config"
	"github.com/supressionstop/xenking_test_1/internal/worker"
	"log/slog"
	"net"
	"net/http"
)

type Http struct {
	cfg        *config.Config
	logger     *slog.Logger
	httpServer *http.Server
	workerPool *worker.Pool
	ErrChan    chan error
}

func NewHttp(cfg *config.Config, logger *slog.Logger, workerPool *worker.Pool) *Http {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	if cfg.App.Environment == "dev" {
		r.Use(middleware.Logger)
	} else {
		r.Use(slogchi.New(logger))
	}
	r.Get("/ready", readyHandler(workerPool))

	httpServer := http.Server{
		Addr:    net.JoinHostPort(cfg.HttpServer.Host, cfg.HttpServer.Port),
		Handler: r,
	}

	return &Http{
		cfg:        cfg,
		logger:     logger,
		httpServer: &httpServer,
		workerPool: workerPool,
	}
}

func (srv *Http) Start() {
	go func() {
		srv.ErrChan <- srv.httpServer.ListenAndServe()
	}()
	srv.logger.Info("http server started", slog.String("addr", srv.httpServer.Addr))
}

func (srv *Http) Shutdown(ctx context.Context) error {
	err := srv.httpServer.Shutdown(ctx)
	if err != nil {
		return err
	}
	srv.logger.Info("http server stopped")
	return nil
}

func readyHandler(workerPool *worker.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		synced := workerPool.IsSynced()
		if synced {
			w.WriteHeader(http.StatusOK)
		} else {
			w.WriteHeader(http.StatusExpectationFailed)
		}
	}
}
