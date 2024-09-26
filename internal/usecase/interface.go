package usecase

import (
	"context"

	"github.com/supressionstop/xenking_test_1/internal/entity"
	"github.com/supressionstop/xenking_test_1/internal/usecase/dto"
	"github.com/supressionstop/xenking_test_1/internal/usecase/enum"
)

// Use cases.
type (
	GetLineFromProvider interface {
		Execute(ctx context.Context, sport string) (entity.Line, error)
	}

	SaveLineToStorage interface {
		Execute(ctx context.Context, line entity.Line) (entity.Line, error)
	}

	// FetchLine from provider and save to storage.
	FetchLine interface {
		Execute(ctx context.Context, sport string) error
	}

	// GetRecentLines from storage.
	GetRecentLines interface {
		Execute(ctx context.Context, sports []enum.Sport) (dto.LineMap, error)
	}

	// CalculateDiff between two collection of lines.
	CalculateDiff interface {
		Execute(prev, curr dto.LineMap) (dto.LinesDiff, error)
	}
)

// Other interfaces.

type LineRepository interface {
	Save(ctx context.Context, line entity.Line) (entity.Line, error)
	GetRecentLines(ctx context.Context, sports []enum.Sport) ([]entity.Line, error)
}

type LineProvider interface {
	GetLine(ctx context.Context, sport string) (dto.ProviderLine, error)
}
