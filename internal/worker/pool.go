package worker

import (
	"context"
	"github.com/supressionstop/xenking_test_1/internal/config"
	"github.com/supressionstop/xenking_test_1/internal/usecase"
	"log/slog"
)

type Pool struct {
	workers []*Worker
	logger  *slog.Logger
}

func NewPool(cfg config.Config, logger *slog.Logger, fetchLine usecase.FetchLine) *Pool {
	pool := &Pool{
		logger:  logger,
		workers: make([]*Worker, len(cfg.Workers)),
	}

	for idx, wc := range cfg.Workers {
		worker := NewWorker(wc.Sport, wc.PollInterval, logger, fetchLine)
		pool.workers[idx] = worker
	}

	return pool
}

func (p *Pool) StartWorkers(ctx context.Context) {
	for _, worker := range p.workers {
		worker.Start(ctx)
	}
}

func (p *Pool) IsSynced() bool {
	for _, worker := range p.workers {
		if !worker.Synced() {
			return false
		}
	}

	return true
}
