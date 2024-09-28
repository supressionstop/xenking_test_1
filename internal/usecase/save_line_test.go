package usecase

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/supressionstop/xenking_test_1/internal/usecase/entity"
	mocks "github.com/supressionstop/xenking_test_1/test/mock"
)

func TestSaveLineUseCaseSuccess(t *testing.T) {
	ctx := context.Background()
	lineToSave := entity.Line{
		ID:          0,
		Name:        entity.Baseball,
		Coefficient: "333.222",
		SavedAt:     time.Time{},
	}
	lineAfter := entity.Line{
		ID:          1,
		Name:        entity.Baseball,
		Coefficient: "333.222",
		SavedAt:     time.Now(),
	}
	lineRepoMock := new(mocks.LineRepository)
	lineRepoMock.On("Save", ctx, lineToSave).Return(lineAfter, nil).Once()
	sut := NewSaveLineUseCase(lineRepoMock)

	savedLine, err := sut.Execute(ctx, lineToSave)

	assert.NoError(t, err)
	assert.Equal(t, lineAfter, savedLine)
}
