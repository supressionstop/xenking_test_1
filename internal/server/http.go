package server

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	slogchi "github.com/samber/slog-chi"
	"log"
	"log/slog"
	"net"
	"net/http"
	"xenking_test_1/internal/config"
)

type Http struct {
	cfg        *config.Config
	logger     *slog.Logger
	httpServer *http.Server
}

func NewHTTP(cfg *config.Config, logger *slog.Logger) *Http {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	if cfg.App.Environment == "dev" {
		r.Use(middleware.Logger)
	} else {
		r.Use(slogchi.New(logger))
	}
	r.Get("/health", healthHandler())
	r.Get("/ready", readyHandler())

	httpServer := http.Server{
		Addr:    net.JoinHostPort(cfg.HttpServer.Host, cfg.HttpServer.Port),
		Handler: r,
	}

	return &Http{
		cfg:        cfg,
		logger:     logger,
		httpServer: &httpServer,
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

func readyHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write([]byte("todo"))
		if err != nil {
			http.Error(w, fmt.Sprint(err), http.StatusInternalServerError)
		}
		w.WriteHeader(http.StatusOK)
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
