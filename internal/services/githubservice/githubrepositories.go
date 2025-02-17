package githubservice

import (
	"context"
	"fmt"

	"github.com/google/go-github/v58/github"
)

type GitHubRepositoriesServiceInterface interface {
	GetRepo(repo string) (*github.Repository, error)
	GetFileContent(repo, path string) (string, error)
}

type GitHubRepositoriesService struct {
	client GitHubRepositoriesInterface
	owner  string
}

func NewGitHubRepositoriesService(client GitHubRepositoriesInterface) *GitHubRepositoriesService {
	return &GitHubRepositoriesService{client: client, owner: "motain"}
}

// Get repository details
func (gh *GitHubRepositoriesService) GetRepo(repo string) (*github.Repository, error) {
	ctx := context.Background()
	repository, _, err := gh.client.Get(ctx, gh.owner, repo)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch repo: %w", err)
	}
	return repository, nil
}

// Get file contents
func (gh *GitHubRepositoriesService) GetFileContent(repo, path string) (string, error) {
	ctx := context.Background()
	fileContent, _, _, err := gh.client.GetContents(ctx, gh.owner, repo, path, nil)
	if err != nil {
		return "", fmt.Errorf("failed to fetch file: %w", err)
	}

	content, err := fileContent.GetContent()
	if err != nil {
		return "", fmt.Errorf("failed to decode file content: %w", err)
	}

	return content, nil
}
