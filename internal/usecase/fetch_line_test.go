package usecase

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/supressionstop/xenking_test_1/internal/usecase/entity"
	mocks "github.com/supressionstop/xenking_test_1/test/mock"
)

func TestFetchLineUseCaseSuccess(t *testing.T) {
	// arrange
	ctx := context.Background()
	sport := entity.Baseball
	line := entity.Line{
		Name:        entity.Baseball,
		Coefficient: "1.5",
	}
	lineSaved := entity.Line{
		ID:          1,
		Name:        entity.Baseball,
		Coefficient: "1.5",
		SavedAt:     time.Now(),
	}
	getLineFromProviderMock := new(mocks.GetLineFromProvider)
	getLineFromProviderMock.On("Execute", ctx, sport).Return(line, nil).Once()
	saveLineToStorageMock := new(mocks.SaveLineToStorage)
	saveLineToStorageMock.On("Execute", ctx, line).Return(lineSaved, nil).Once()
	sut := NewFetchLineUseCase(getLineFromProviderMock, saveLineToStorageMock)

	// act
	err := sut.Execute(ctx, sport)

	// assert
	assert.NoError(t, err)
	getLineFromProviderMock.AssertExpectations(t)
	saveLineToStorageMock.AssertExpectations(t)
}
