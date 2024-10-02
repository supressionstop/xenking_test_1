package usecase

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/supressionstop/xenking_test_1/internal/usecase/entity"
	mocks "github.com/supressionstop/xenking_test_1/test/mock"
)

func TestGetLineSuccess(t *testing.T) {
	// arrange
	ctx := context.Background()
	sport := entity.Baseball
	providerLine := entity.Line{
		Name:        entity.Baseball,
		Coefficient: "1.2345",
	}
	want := entity.Line{
		ID:          0,
		Name:        entity.Baseball,
		Coefficient: "1.2345",
		SavedAt:     time.Time{},
	}
	providerMock := new(mocks.LineProvider)
	providerMock.On("GetLine", ctx, sport).Return(providerLine, nil).Once()
	sut := NewGetLineUseCase(providerMock)

	// act
	line, err := sut.Execute(ctx, sport)

	// assert
	assert.NoError(t, err)
	assert.Equal(t, want, line)
	providerMock.AssertExpectations(t)
}
