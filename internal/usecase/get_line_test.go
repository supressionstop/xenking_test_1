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

func TestGetLineSuccess(t *testing.T) {
	// arrange
	ctx := context.Background()
	sport := enum.Baseball
	providerLine := dto.ProviderLine{
		Sport: enum.Baseball,
		Rate:  "1.2345",
	}
	want := entity.Line{
		ID:          0,
		Name:        enum.Baseball,
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
