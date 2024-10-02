package usecase

import (
	"context"

	"github.com/supressionstop/xenking_test_1/internal/usecase/entity"
)

// RecentLinesGetter from storage.
type RecentLinesGetter interface {
	Execute(ctx context.Context, sports []entity.Sport) (entity.LineMap, error)
}

type GetRecentLines struct {
	lineRepository LineRepository
}

func NewGetRecentLinesUseCase(lineRepository LineRepository) *GetRecentLines {
	return &GetRecentLines{
		lineRepository: lineRepository,
	}
}

func (uc *GetRecentLines) Execute(ctx context.Context, sports []entity.Sport) (entity.LineMap, error) {
	lines, err := uc.lineRepository.GetRecentLines(ctx, sports)
	if err != nil {
		return nil, err
	}

	return entity.LineMapFromSports(lines), nil
}
