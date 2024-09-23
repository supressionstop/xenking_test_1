package usecase

import (
	"context"
	"xenking_test_1/internal/entity"
	"xenking_test_1/internal/usecase/policy"
)

// GetLineUseCase gets line from provider
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
