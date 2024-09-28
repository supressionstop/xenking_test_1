package usecase

import (
	"context"

	"github.com/supressionstop/xenking_test_1/internal/usecase/entity"
)

type SaveLineToStorage interface {
	Execute(ctx context.Context, line entity.Line) (entity.Line, error)
}

type LineRepository interface {
	Save(ctx context.Context, line entity.Line) (entity.Line, error)
	GetRecentLines(ctx context.Context, sports []entity.Sport) ([]entity.Line, error)
}

type SaveLine struct {
	lineRepo LineRepository
}

func NewSaveLineUseCase(lineRepository LineRepository) *SaveLine {
	return &SaveLine{
		lineRepo: lineRepository,
	}
}

func (s *SaveLine) Execute(ctx context.Context, line entity.Line) (entity.Line, error) {
	return s.lineRepo.Save(ctx, line)
}
