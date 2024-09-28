package worker

import (
	"context"
	"fmt"
	"log/slog"
	"sync"

	"github.com/supressionstop/xenking_test_1/internal/application/task"
)

type Pool struct {
	workerCount        int
	tasks              <-chan *task.Task
	logger             *slog.Logger
	wg                 *sync.WaitGroup
	workerSync         chan struct{}
	AllWorkersSynced   chan struct{}
	isAllWorkersSynced bool
}

func NewPool(workerCount int, tasks <-chan *task.Task, logger *slog.Logger) *Pool {
	return &Pool{
		workerCount:      workerCount,
		tasks:            tasks,
		logger:           logger,
		wg:               &sync.WaitGroup{},
		workerSync:       make(chan struct{}, workerCount),
		AllWorkersSynced: make(chan struct{}),
	}
}

func (p *Pool) Start(ctx context.Context) {
	for n := 1; n <= p.workerCount; n++ {
		worker := NewWorker(p.tasks, fmt.Sprintf("worker-%d", n), p.logger, p.workerSync)
		worker.Start(ctx, p.wg)
	}

	go func() {
		left := p.workerCount
		for range p.workerSync {
			left--
			if left == 0 {
				p.isAllWorkersSynced = true
				p.AllWorkersSynced <- struct{}{}
			}
		}
	}()

	p.wg.Wait()
}

func (p *Pool) IsSynced() bool {
	return p.isAllWorkersSynced
}
