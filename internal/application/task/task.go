package task

import (
	"context"

	"github.com/supressionstop/xenking_test_1/internal/usecase"
	"github.com/supressionstop/xenking_test_1/internal/usecase/entity"
)

type Task struct {
	sport     entity.Sport
	fetchLine usecase.FetchLiner
}

func NewTask(sport entity.Sport, fetchLine usecase.FetchLiner) *Task {
	return &Task{
		sport:     sport,
		fetchLine: fetchLine,
	}
}

func (t *Task) Process(ctx context.Context) error {
	return t.fetchLine.Execute(ctx, t.sport)
}

func (t *Task) String() string {
	return t.sport
}
