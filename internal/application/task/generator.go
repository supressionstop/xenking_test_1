package task

import (
	"context"
	"time"

	"github.com/supressionstop/xenking_test_1/internal/infrastructure/config"
	"github.com/supressionstop/xenking_test_1/internal/usecase"
)

type Generator struct {
	settings  []config.Worker
	tasks     chan *Task
	fetchLine usecase.FetchLiner
}

func NewGenerator(settings []config.Worker, fetchLine usecase.FetchLiner) *Generator {
	tasks := make(chan *Task, len(settings))
	return &Generator{
		settings:  settings,
		tasks:     tasks,
		fetchLine: fetchLine,
	}
}

func (tg *Generator) Start(ctx context.Context) {
	for _, setting := range tg.settings {
		go func() {
			ticker := time.NewTicker(setting.PollInterval)
			defer ticker.Stop()
			for {
				select {
				case <-ticker.C:
					tg.tasks <- NewTask(setting.Sport, tg.fetchLine)
				case <-ctx.Done():
					return
				}
			}
		}()
	}
}

func (tg *Generator) Tasks() <-chan *Task {
	return tg.tasks
}
