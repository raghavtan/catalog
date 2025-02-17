//go:build unit
// +build unit

package repository_test

import (
	"testing"

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

func TestFetchData(t *testing.T) {
	cfg := &MockConfigService{}
	repo := repository.NewRepository(cfg)

	result := repo.FetchData()
	assert.Contains(t, result, "Data")
}
