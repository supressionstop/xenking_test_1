package usecase

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/supressionstop/xenking_test_1/internal/entity"
	"github.com/supressionstop/xenking_test_1/internal/usecase/dto"
	"github.com/supressionstop/xenking_test_1/internal/usecase/enum"
	mocks "github.com/supressionstop/xenking_test_1/test/mock"
	"testing"
	"time"
)

func TestGetRecentLineUseCaseSuccess(t *testing.T) {
	ctx := context.Background()
	sports := []string{enum.Baseball, enum.Football, enum.Soccer}
	entities := []entity.Line{
		{
			ID:          1,
			Name:        enum.Baseball,
			Coefficient: "1.1111",
			SavedAt:     time.Now().AddDate(0, 0, -1),
		},
		{
			ID:          2,
			Name:        enum.Football,
			Coefficient: "2.2222",
			SavedAt:     time.Now().AddDate(0, 0, -2),
		},
		{
			ID:          3,
			Name:        enum.Soccer,
			Coefficient: "3.3333",
			SavedAt:     time.Now().AddDate(0, 0, -3),
		},
	}
	want := dto.LineMap{
		enum.Baseball: entities[0],
		enum.Football: entities[1],
		enum.Soccer:   entities[2],
	}
	lineRepoMock := new(mocks.LineRepository)
	lineRepoMock.On("GetRecentLines", ctx, sports).Return(entities, nil).Once()
	sut := NewGetRecentLinesUseCase(lineRepoMock)

	recentLines, err := sut.Execute(ctx, sports)

	assert.NoError(t, err)
	assert.Len(t, recentLines, len(want))
	assert.Equal(t, want, recentLines)
}
