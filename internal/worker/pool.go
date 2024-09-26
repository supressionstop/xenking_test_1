package worker

import (
	"context"
	"log"
	"log/slog"

	"github.com/supressionstop/xenking_test_1/internal/config"
	"github.com/supressionstop/xenking_test_1/internal/usecase"
)

type Pool struct {
	workers        []*Worker
	logger         *slog.Logger
	FirstSyncChan  chan struct{}
	workerSyncChan chan struct{}
	workersErrChan chan error
}

func NewPool(cfg config.Config, logger *slog.Logger, fetchLine usecase.FetchLine) *Pool {
	pool := &Pool{
		logger:         logger,
		workers:        make([]*Worker, len(cfg.Workers)),
		FirstSyncChan:  make(chan struct{}),
		workerSyncChan: make(chan struct{}, len(cfg.Workers)),
		workersErrChan: make(chan error, len(cfg.Workers)),
	}

	for idx, wc := range cfg.Workers {
		worker := NewWorker(wc.Sport, wc.PollInterval, logger, fetchLine, pool.workerSyncChan, pool.workersErrChan)
		pool.workers[idx] = worker
	}

	return pool
}

func (p *Pool) StartWorkers(ctx context.Context) {
	p.logger.Info("starting workers...")
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
			case err := <-p.workersErrChan:
				p.logger.Error("worker err:", slog.Any("err", err))
				log.Fatal(err)
			case <-ctx.Done():
				return
			}
		}
	}()
}
