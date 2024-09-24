package usecase

import (
	"context"
	"github.com/supressionstop/xenking_test_1/internal/entity"
	"github.com/supressionstop/xenking_test_1/internal/usecase/dto"
	"github.com/supressionstop/xenking_test_1/internal/usecase/enum"
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

	GetRecentLines interface {
		Execute(ctx context.Context, sports []enum.Sport) (dto.LineMap, error)
	}

	CalculateDiff interface {
		Execute(prev, curr dto.LineMap) (dto.LinesDiff, error)
	}
)

type LineRepository interface {
	Save(context.Context, entity.Line) (entity.Line, error)
	GetRecentLines(ctx context.Context, sports []enum.Sport) ([]entity.Line, error)
}

type LineProvider interface {
	GetLine(ctx context.Context, sport string) (dto.ProviderLine, error)
}
