//go:build unit
// +build unit

package githubservice_test

import (
	"context"
	"errors"
	"testing"

	"github.com/google/go-github/v58/github"
	"github.com/motain/fact-collector/internal/services/githubservice"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockGitHubClient is a mock of the GitHub client
type MockGitHubClient struct {
	mock.Mock
}

// Ensure MockGitHubClient implements GitHubRepositoriesInterface
var _ githubservice.GitHubRepositoriesInterface = (*MockGitHubClient)(nil)

// Get is a mock implementation of GitHubRepositoriesInterface.Get
func (m *MockGitHubClient) Get(ctx context.Context, owner, repo string) (*github.Repository, *github.Response, error) {
	args := m.Called(ctx, owner, repo)
	return args.Get(0).(*github.Repository), args.Get(1).(*github.Response), args.Error(2)
}

// GetContents is a mock implementation of GitHubRepositoriesInterface.GetContents
func (m *MockGitHubClient) GetContents(ctx context.Context, owner, repo, path string, opts *github.RepositoryContentGetOptions) (fileContent *github.RepositoryContent, directoryContent []*github.RepositoryContent, resp *github.Response, err error) {
	args := m.Called(ctx, owner, repo, path, opts)
	return args.Get(0).(*github.RepositoryContent), args.Get(1).([]*github.RepositoryContent), nil, args.Error(2)
}

func TestGetFileContent(t *testing.T) {
	mockClient := new(MockGitHubClient)
	service := githubservice.NewGitHubRepositoriesService(mockClient)

	t.Run("successful content retrieval", func(t *testing.T) {
		expectedContent := "file content"
		mockFileContent := &github.RepositoryContent{Content: github.String(expectedContent)}
		mockClient.On("GetContents", mock.Anything, "motain", "repo", "path", (*github.RepositoryContentGetOptions)(nil)).Return(mockFileContent, ([]*github.RepositoryContent)(nil), nil).Once()

		content, err := service.GetFileContent("repo", "path")
		assert.NoError(t, err)
		assert.Equal(t, expectedContent, content)
		mockClient.AssertExpectations(t)
	})

	t.Run("error fetching file", func(t *testing.T) {
		mockClient.On("GetContents", mock.Anything, "motain", "repo", "path", (*github.RepositoryContentGetOptions)(nil)).Return((*github.RepositoryContent)(nil), ([]*github.RepositoryContent)(nil), errors.New("failed to fetch file")).Once()

		content, err := service.GetFileContent("repo", "path")
		assert.Error(t, err)
		assert.Equal(t, "", content)
		mockClient.AssertExpectations(t)
	})

	t.Run("error decoding file content", func(t *testing.T) {
		encoding := "none"
		mockFileContent := &github.RepositoryContent{Content: github.String(""), Encoding: &encoding}
		mockClient.On("GetContents", mock.Anything, "motain", "repo", "path", (*github.RepositoryContentGetOptions)(nil)).Return(mockFileContent, ([]*github.RepositoryContent)(nil), nil).Once()

		content, err := service.GetFileContent("repo", "path")
		assert.Error(t, err)
		assert.Equal(t, "", content)
		mockClient.AssertExpectations(t)
	})
}
