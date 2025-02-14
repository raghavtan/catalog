package githubservice

import (
	"context"
	"fmt"

	"github.com/google/go-github/v58/github"
	"github.com/motain/fact-collector/internal/services/configservice"
	"github.com/motain/fact-collector/internal/services/keyringservice"
	"golang.org/x/oauth2"
)

type GitHubServiceInterface interface {
	GetRepo(owner, repo string) (*github.Repository, error)
	GetFileContent(owner, repo, path string) (string, error)
}

type GitHubService struct {
	client *github.Client
}

func NewGitHubService(
	cfg configservice.ConfigServiceInterface,
	kr keyringservice.KeyringServiceInterface,
) *GitHubService {
	serviceName := "gh:github.com"
	token, tokenErr := kr.Get(serviceName, cfg.GetGithubUser())
	if tokenErr != nil {
		panic(tokenErr)
	}

	ctx := context.Background()
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})
	tc := oauth2.NewClient(ctx, ts)

	return &GitHubService{client: github.NewClient(tc)}
}

// Get repository details
func (gh *GitHubService) GetRepo(owner, repo string) (*github.Repository, error) {
	ctx := context.Background()
	repository, _, err := gh.client.Repositories.Get(ctx, owner, repo)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch repo: %w", err)
	}
	return repository, nil
}

// Get file contents
func (gh *GitHubService) GetFileContent(owner, repo, path string) (string, error) {
	ctx := context.Background()
	fileContent, _, _, err := gh.client.Repositories.GetContents(ctx, owner, repo, path, nil)
	if err != nil {
		return "", fmt.Errorf("failed to fetch file: %w", err)
	}

	content, err := fileContent.GetContent()
	if err != nil {
		return "", fmt.Errorf("failed to decode file content: %w", err)
	}

	return content, nil
}
