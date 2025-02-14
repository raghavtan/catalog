//go:build unit
// +build unit

package repository_test

import (
	"testing"

	"github.com/google/go-github/v58/github"
	"github.com/motain/fact-collector/internal/modules/sample/repository"
	"github.com/stretchr/testify/assert"
)

type MockConfigService struct{}

func (m *MockConfigService) GetGithubOrg() string {
	return "Mocked Data"
}

func (m *MockConfigService) GetGithubToken() string {
	return "Mocked Token"
}

func (m *MockConfigService) GetGithubUser() string {
	return "Mocked User"
}

type MockGitHubService struct{}

func (m *MockGitHubService) GetFileContent(owner, repo, path string) (string, error) {
	return "Mocked File Content", nil
}

func (m *MockGitHubService) GetRepo(owner, repo string) (*github.Repository, error) {
	name := "Mocked Repo"
	return &github.Repository{Name: &name}, nil
}

func TestFetchData(t *testing.T) {
	cfg := &MockConfigService{}
	gh := &MockGitHubService{}
	repo := repository.NewRepository(cfg, gh)

	result := repo.FetchData()
	assert.Contains(t, result, "Data from GitHub Org: ")
}
