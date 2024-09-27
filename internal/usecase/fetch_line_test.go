package usecase

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/supressionstop/xenking_test_1/internal/entity"
	"github.com/supressionstop/xenking_test_1/internal/usecase/enum"
	mocks "github.com/supressionstop/xenking_test_1/test/mock"
	"testing"
	"time"
)

func TestFetchLineUseCaseSuccess(t *testing.T) {
	// arrange
	ctx := context.Background()
	sport := enum.Baseball
	line := entity.Line{
		Name:        enum.Baseball,
		Coefficient: "1.5",
	}
	lineSaved := entity.Line{
		ID:          1,
		Name:        enum.Baseball,
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
