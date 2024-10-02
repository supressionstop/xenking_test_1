package usecase

import (
	"context"

	"github.com/supressionstop/xenking_test_1/internal/usecase/entity"
)

// GetLiner - gets line from provider.
type GetLiner interface {
	Execute(ctx context.Context, sport string) (entity.Line, error)
}

type LineProvider interface {
	GetLine(ctx context.Context, sport string) (entity.Line, error)
}

type GetLine struct {
	provider LineProvider
}

func NewGetLineUseCase(provider LineProvider) *GetLine {
	return &GetLine{
		provider: provider,
	}
}

func (uc *GetLine) Execute(ctx context.Context, sportName string) (entity.Line, error) {
	line, err := uc.provider.GetLine(ctx, sportName)
	if err != nil {
		return entity.Line{}, err
	}

	return line, nil
}
