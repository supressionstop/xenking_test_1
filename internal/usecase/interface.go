package usecase

import (
	"context"
	"xenking_test_1/internal/entity"
	"xenking_test_1/internal/usecase/dto"
)

type (
	GetLine interface {
		Execute(ctx context.Context, sport string) (entity.Line, error)
	}

	SaveLine interface {
		Execute(ctx context.Context, line entity.Line) (entity.Line, error)
	}

	FetchLine interface {
		Execute(ctx context.Context, sport string) error
	}
)

type LineRepository interface {
	Save(context.Context, entity.Line) (entity.Line, error)
}

type LineProvider interface {
	GetLine(ctx context.Context, sport string) (dto.ProviderLine, error)
}
