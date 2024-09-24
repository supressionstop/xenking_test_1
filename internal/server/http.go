package server

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	slogchi "github.com/samber/slog-chi"
	"github.com/supressionstop/xenking_test_1/internal/config"
	"github.com/supressionstop/xenking_test_1/internal/worker"
	"log"
	"log/slog"
	"net"
	"net/http"
)

type Http struct {
	cfg        *config.Config
	logger     *slog.Logger
	httpServer *http.Server
	workerPool *worker.Pool
}

func NewHTTP(cfg *config.Config, logger *slog.Logger, workerPool *worker.Pool) *Http {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	if cfg.App.Environment == "dev" {
		r.Use(middleware.Logger)
	} else {
		r.Use(slogchi.New(logger))
	}
	r.Get("/health", healthHandler())
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

func (srv *Http) Start(ctx context.Context) {
	// TODO: use context
	go func() {
		err := srv.httpServer.ListenAndServe()
		if err != nil {
			log.Fatal(fmt.Errorf("http server error (TODO REMOVE): %w", err))
		}
	}()
	srv.logger.Info("http server started", slog.String("addr", srv.httpServer.Addr))
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

func healthHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		type Health struct {
			Ok      bool   `json:"ok"`
			Message string `json:"message"`
		}
		body := Health{
			Ok:      true,
			Message: "http server is working",
		}
		payload, err := json.Marshal(body)
		if err != nil {
			log.Fatal(fmt.Errorf("failed to marshal health response TODO REMOVE: %w", err))
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, err = w.Write(payload)
		if err != nil {
			log.Fatal(fmt.Errorf("failed to write health response TODO REMOVE: %w", err))
		}
	}
}
