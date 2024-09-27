package worker

import (
	"context"
	"log/slog"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/supressionstop/xenking_test_1/internal/usecase"
)

type Worker struct {
	fetchLine         usecase.FetchLine
	sport             string
	id                string
	interval          time.Duration
	ErrCh             chan error
	isFirstTimeSynced bool
	logger            *slog.Logger
	mu                sync.RWMutex
	workerSyncChan    chan struct{}
}

func NewWorker(
	sport string,
	interval time.Duration,
	logger *slog.Logger,
	fetchLine usecase.FetchLine,
	workerSyncChan chan struct{},
	errCh chan error,
) *Worker {
	return &Worker{
		id:                uuid.New().String(),
		fetchLine:         fetchLine,
		sport:             sport,
		interval:          interval,
		logger:            logger,
		workerSyncChan:    workerSyncChan,
		ErrCh:             errCh,
		isFirstTimeSynced: false,
		mu:                sync.RWMutex{},
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
					w.ErrCh <- err
					continue
				} // TODO: handle errors

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
	w.workerSyncChan <- struct{}{}
}
