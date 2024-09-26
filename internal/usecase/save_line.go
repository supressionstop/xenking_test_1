package usecase

import (
	"context"

	"github.com/supressionstop/xenking_test_1/internal/entity"
)

type SaveLineUseCase struct {
	lineRepo LineRepository
}

func NewSaveLineUseCase(lineRepository LineRepository) *SaveLineUseCase {
	return &SaveLineUseCase{
		lineRepo: lineRepository,
	}
}

func (s *SaveLineUseCase) Execute(ctx context.Context, line entity.Line) (entity.Line, error) {
	return s.lineRepo.Save(ctx, line)
}
