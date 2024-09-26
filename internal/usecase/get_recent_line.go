package usecase

import (
	"context"

	"github.com/supressionstop/xenking_test_1/internal/usecase/dto"
	"github.com/supressionstop/xenking_test_1/internal/usecase/enum"
)

type GetRecentLinesUseCase struct {
	lineRepository LineRepository
}

func NewGetRecentLinesUseCase(lineRepository LineRepository) *GetRecentLinesUseCase {
	return &GetRecentLinesUseCase{
		lineRepository: lineRepository,
	}
}

func (uc *GetRecentLinesUseCase) Execute(ctx context.Context, sports []enum.Sport) (dto.LineMap, error) {
	lines, err := uc.lineRepository.GetRecentLines(ctx, sports)
	if err != nil {
		return nil, err
	}

	return dto.LineMapFromSports(lines), nil
}
