package worker

import (
	"context"
	"github.com/supressionstop/xenking_test_1/internal/config"
	"github.com/supressionstop/xenking_test_1/internal/usecase"
	"log/slog"
)

type Pool struct {
	workers        []*Worker
	logger         *slog.Logger
	FirstSyncChan  chan struct{}
	workerSyncChan chan struct{}
}

func NewPool(cfg config.Config, logger *slog.Logger, fetchLine usecase.FetchLine) *Pool {
	pool := &Pool{
		logger:         logger,
		workers:        make([]*Worker, len(cfg.Workers)),
		FirstSyncChan:  make(chan struct{}),
		workerSyncChan: make(chan struct{}, len(cfg.Workers)),
	}

	for idx, wc := range cfg.Workers {
		worker := NewWorker(wc.Sport, wc.PollInterval, logger, fetchLine, pool.workerSyncChan)
		pool.workers[idx] = worker
	}

	return pool
}

func (p *Pool) StartWorkers(ctx context.Context) {
	for _, worker := range p.workers {
		worker.Start(ctx)
	}

	p.waitSynchronization(ctx)
}

func (p *Pool) IsSynced() bool {
	for _, worker := range p.workers {
		if !worker.Synced() {
			return false
		}
	}

	return true
}

func (p *Pool) waitSynchronization(ctx context.Context) {
	go func() {
		leftToSync := len(p.workers)
		for {
			select {
			case <-p.workerSyncChan:
				leftToSync--
				if leftToSync == 0 {
					p.FirstSyncChan <- struct{}{}
					close(p.FirstSyncChan)
					return
				}
			case <-ctx.Done():
				return
			}
		}
	}()
}
