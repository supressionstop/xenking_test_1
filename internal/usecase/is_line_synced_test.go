package usecase

import (
	"testing"

	"github.com/stretchr/testify/assert"

	mocks "github.com/supressionstop/xenking_test_1/test/mock"
)

func TestAppReady(t *testing.T) {
	checker := new(mocks.LineChecker)
	checker.On("IsSynced").Return(true)
	sut := NewIsLineSynced(checker)

	assert.True(t, sut.Execute())

	checker.AssertExpectations(t)
}

func TestAppIsNotReady(t *testing.T) {
	checker := new(mocks.LineChecker)
	checker.On("IsSynced").Return(false)
	sut := NewIsLineSynced(checker)

	assert.False(t, sut.Execute())

	checker.AssertExpectations(t)
}
