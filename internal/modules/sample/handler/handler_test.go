//go:build unit
// +build unit

package handler_test

import (
	"testing"

	"github.com/motain/fact-collector/internal/modules/sample/handler"
	"github.com/stretchr/testify/assert"
)

type MockRepo struct{}

func (m *MockRepo) FetchData() string {
	return "Mocked Data"
}

func TestHandle(t *testing.T) {
	repo := &MockRepo{}
	h := handler.NewHandler(repo)

	result := h.Handle()
	assert.Equal(t, "Mocked Data", result)
}
