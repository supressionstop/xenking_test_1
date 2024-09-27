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

func TestSaveLineUseCaseSuccess(t *testing.T) {
	ctx := context.Background()
	lineToSave := entity.Line{
		ID:          0,
		Name:        enum.Baseball,
		Coefficient: "333.222",
		SavedAt:     time.Time{},
	}
	lineAfter := entity.Line{
		ID:          1,
		Name:        enum.Baseball,
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
