package server

import (
	"context"
	"errors"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	slogchi "github.com/samber/slog-chi"

	controller "github.com/supressionstop/xenking_test_1/internal/application/controller/http"
	"github.com/supressionstop/xenking_test_1/internal/infrastructure/config"
)

type Http struct {
	cfg        *config.Config
	logger     *slog.Logger
	httpServer *http.Server
	Finish     chan struct{}
}

func NewHttp(cfg *config.Config, logger *slog.Logger, readyController *controller.ReadyController) *Http {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(slogchi.New(logger))
	r.Get("/ready", readyController.Ready)

	httpServer := &http.Server{
		Addr:    ":" + cfg.HTTPServer.Port,
		Handler: r,
	}

	return &Http{
		cfg:        cfg,
		logger:     logger,
		httpServer: httpServer,
		Finish:     make(chan struct{}, 1),
	}
}

func (srv *Http) Start(ctx context.Context) {
	go func() {
		<-ctx.Done()
		err := srv.httpServer.Shutdown(ctx)
		if err != nil {
			srv.logger.Error("HTTP server stop error", slog.Any("error", err))
			return
		}
		srv.logger.Info("HTTP server stopped")
		srv.Finish <- struct{}{}
	}()

	srv.logger.Info("HTTP server started", slog.String("addr", srv.httpServer.Addr))
	if err := srv.httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		srv.logger.Error("HTTP server start error", slog.Any("error", err))
		panic(err)
	}
}
