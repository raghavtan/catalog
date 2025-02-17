//go:build unit
// +build unit

package handler_test

import (
	"testing"

	"github.com/google/go-github/v58/github"
	"github.com/motain/fact-collector/internal/modules/sample/handler"
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

func TestHandle(t *testing.T) {
	gh := &MockGitHubRepositoriesService{}

	h := handler.NewHandler(gh)

	result := h.Handle()
	assert.Equal(t, "", result)
}
