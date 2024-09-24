package usecase

import (
	"github.com/shopspring/decimal"
	"github.com/supressionstop/xenking_test_1/internal/usecase/dto"
)

type CalculateDiffUseCase struct {
}

func NewCalculateDiffUseCase() *CalculateDiffUseCase {
	return &CalculateDiffUseCase{}
}

func (uc *CalculateDiffUseCase) Execute(prev, curr dto.LineMap) (dto.LinesDiff, error) {
	noPrev := len(prev) == 0
	if noPrev {
		return uc.firstResponse(curr)
	}

	return uc.diff(prev, curr)
}

func (uc *CalculateDiffUseCase) firstResponse(curr dto.LineMap) (dto.LinesDiff, error) {
	firstResponse := make(dto.LinesDiff, len(curr))
	for sport, line := range curr {
		firstResponse[sport] = line.Coefficient
	}
	return firstResponse, nil
}

func (uc *CalculateDiffUseCase) diff(prev, curr dto.LineMap) (dto.LinesDiff, error) {
	diff := make(dto.LinesDiff, len(curr))
	for sport, line := range prev {
		prevRate, err := decimal.NewFromString(line.Coefficient)
		if err != nil {
			return nil, err
		}
		currRate, err := decimal.NewFromString(curr[sport].Coefficient)
		if err != nil {
			return nil, err
		}
		diff[sport] = prevRate.Sub(currRate).String()
	}
	return diff, nil
}
