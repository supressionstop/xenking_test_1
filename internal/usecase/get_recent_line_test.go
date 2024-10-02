package usecase

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/supressionstop/xenking_test_1/internal/usecase/entity"
	mocks "github.com/supressionstop/xenking_test_1/test/mock"
)

func TestGetRecentLineUseCaseSuccess(t *testing.T) {
	ctx := context.Background()
	sports := []string{entity.Baseball, entity.Football, entity.Soccer}
	entities := []entity.Line{
		{
			ID:          1,
			Name:        entity.Baseball,
			Coefficient: "1.1111",
			SavedAt:     time.Now().AddDate(0, 0, -1),
		},
		{
			ID:          2,
			Name:        entity.Football,
			Coefficient: "2.2222",
			SavedAt:     time.Now().AddDate(0, 0, -2),
		},
		{
			ID:          3,
			Name:        entity.Soccer,
			Coefficient: "3.3333",
			SavedAt:     time.Now().AddDate(0, 0, -3),
		},
	}
	want := entity.LineMap{
		entity.Baseball: entities[0],
		entity.Football: entities[1],
		entity.Soccer:   entities[2],
	}
	lineRepoMock := new(mocks.LineRepository)
	lineRepoMock.On("GetRecentLines", ctx, sports).Return(entities, nil).Once()
	sut := NewGetRecentLinesUseCase(lineRepoMock)

	recentLines, err := sut.Execute(ctx, sports)

	assert.NoError(t, err)
	assert.Len(t, recentLines, len(want))
	assert.Equal(t, want, recentLines)
}
