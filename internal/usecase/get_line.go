package usecase

import (
	"context"

	"github.com/supressionstop/xenking_test_1/internal/entity"
	"github.com/supressionstop/xenking_test_1/internal/usecase/policy"
)

type GetLineUseCase struct {
	provider LineProvider
}

func NewGetLineUseCase(provider LineProvider) *GetLineUseCase {
	return &GetLineUseCase{
		provider: provider,
	}
}

func (uc *GetLineUseCase) Execute(ctx context.Context, sportName string) (entity.Line, error) {
	line, err := uc.provider.GetLine(ctx, sportName)
	if err != nil {
		return entity.Line{}, err
	}

	return policy.ProviderLineToEntity(line), nil
}
