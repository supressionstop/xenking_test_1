package worker

import (
	"context"
	"log/slog"
	"sync"

	"github.com/supressionstop/xenking_test_1/internal/application/task"
)

type Worker struct {
	ID         string
	tasks      <-chan *task.Task
	synced     chan struct{}
	logger     *slog.Logger
	isSynced   bool
	errorCount int
	maxErrors  int
}

const maxWorkerErrorDefault = 5

func NewWorker(tasks <-chan *task.Task, ID string, logger *slog.Logger, synced chan struct{}) *Worker {
	return &Worker{
		ID:        ID,
		tasks:     tasks,
		synced:    synced,
		logger:    logger,
		maxErrors: maxWorkerErrorDefault,
	}
}

func (w *Worker) Start(ctx context.Context, wg *sync.WaitGroup) {
	w.logger.Info("starting worker", slog.String("id", w.ID))

	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			select {
			case task := <-w.tasks:
				err := task.Process(ctx)
				if err != nil {
					w.logger.Error(
						"failed to process task",
						slog.String("worker_id", w.ID),
						slog.String("sport", task.String()),
						slog.String("error", err.Error()),
					)
					w.errorCount++
					if w.errorCount >= w.maxErrors {
						w.logger.Error("max errors exceeded, stops worker", slog.String("worker_id", w.ID))
						return
					}
					continue
				}
				w.logger.Debug(
					"line fetched",
					slog.String("worker_id", w.ID),
					slog.String("sport", task.String()),
				)
				w.setSynced()
			case <-ctx.Done():
				w.logger.Info("worker stopped", slog.String("id", w.ID))
				return
			}
		}
	}()
}

func (w *Worker) setSynced() {
	if w.isSynced {
		return
	}
	w.isSynced = true
	w.synced <- struct{}{}
}
