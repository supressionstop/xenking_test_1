package usecase

import (
	"github.com/shopspring/decimal"

	"github.com/supressionstop/xenking_test_1/internal/usecase/entity"
)

// CalculateDiffer between two collection of lines.
type CalculateDiffer interface {
	Execute(prev, curr entity.LineMap) (entity.LinesDiff, error)
}

type CalculateDiff struct{}

func NewCalculateDiffUseCase() *CalculateDiff {
	return &CalculateDiff{}
}

func (uc *CalculateDiff) Execute(prev, curr entity.LineMap) (entity.LinesDiff, error) {
	noPrev := len(prev) == 0
	if noPrev {
		return uc.firstResponse(curr)
	}

	return uc.diff(prev, curr)
}

func (uc *CalculateDiff) firstResponse(curr entity.LineMap) (entity.LinesDiff, error) {
	firstResponse := make(entity.LinesDiff, len(curr))
	for sport, line := range curr {
		firstResponse[sport] = line.Coefficient
	}
	return firstResponse, nil
}

func (uc *CalculateDiff) diff(prev, curr entity.LineMap) (entity.LinesDiff, error) {
	diff := make(entity.LinesDiff, len(curr))
	for sport, line := range prev {
		prevRate, err := decimal.NewFromString(line.Coefficient)
		if err != nil {
			return nil, err
		}
		currRate, err := decimal.NewFromString(curr[sport].Coefficient)
		if err != nil {
			return nil, err
		}
		diff[sport] = currRate.Sub(prevRate).String()
	}
	return diff, nil
}
