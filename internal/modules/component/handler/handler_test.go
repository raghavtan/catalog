//go:build unit
// +build unit

package handler_test

import (
	"context"
	"testing"

	"github.com/google/go-github/v58/github"
	"github.com/motain/fact-collector/internal/modules/component/handler"
	"github.com/stretchr/testify/assert"
)

type MockGitHubRepositoriesService struct{}

func (m *MockGitHubRepositoriesService) GetFileContent(repo, path string) (string, error) {
	return "Mocked File Content", nil
}

func (m *MockGitHubRepositoriesService) GetRepo(repo string) (*github.Repository, error) {
	name := "Mocked Repo"
	return &github.Repository{Name: &name}, nil
}

type MockCompassService struct{}

func (m *MockCompassService) Run(ctx context.Context, query string, variables map[string]interface{}, response interface{}) error {
	return nil
}

func TestHandle(t *testing.T) {
	gh := &MockGitHubRepositoriesService{}
	compass := &MockCompassService{}

	h := handler.NewHandler(gh, compass)

	result := h.Handle()
	assert.Equal(t, "", result)
}
