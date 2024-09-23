package worker

import (
	"context"
	"github.com/google/uuid"
	"log/slog"
	"sync"
	"time"
	"xenking_test_1/internal/usecase"
)

type Worker struct {
	fetchLine         usecase.FetchLine
	sport             string
	id                string
	interval          time.Duration
	errCh             chan error
	isFirstTimeSynced bool
	logger            *slog.Logger
	mu                sync.RWMutex
}

func NewWorker(sport string, interval time.Duration, logger *slog.Logger, fetchLine usecase.FetchLine) *Worker {
	return &Worker{
		id:        uuid.New().String(),
		fetchLine: fetchLine,
		sport:     sport,
		interval:  interval,
		logger:    logger,
	}
}

func (w *Worker) Start(ctx context.Context) {
	ticker := time.NewTicker(w.interval)

	go func() {
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				err := w.fetchLine.Execute(ctx, w.sport)
				if err != nil {
					w.errCh <- err
				}
				w.setSynced()
				w.logger.Debug("line fetched", slog.String("id", w.id), slog.String("sport", w.sport))
			case <-ctx.Done():
				return
			}
		}
	}()

	w.logger.Info("worker started", slog.String("id", w.id), slog.String("sport", w.sport))
}

func (w *Worker) Synced() bool {
	w.mu.RLock()
	defer w.mu.RUnlock()
	return w.isFirstTimeSynced
}

func (w *Worker) setSynced() {
	if w.isFirstTimeSynced {
		return
	}
	w.mu.Lock()
	defer w.mu.Unlock()
	w.isFirstTimeSynced = true
}
