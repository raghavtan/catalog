package githubservice

//go:generate mockgen -destination=./mocks/mock_github_client.go -package=githubservice github.com/motain/of-catalog/internal/services/githubservice GitHubRepositoriesInterface

import (
	"context"
	"fmt"
	"io"

	"github.com/google/go-github/v58/github"
	"github.com/motain/of-catalog/internal/services/configservice"
	"github.com/motain/of-catalog/internal/services/keyringservice"
	"golang.org/x/oauth2"
)

type GitHubRepositoriesInterface interface {
	Get(ctx context.Context, owner, repo string) (*github.Repository, *github.Response, error)
	GetContents(ctx context.Context, owner, repo, path string, opts *github.RepositoryContentGetOptions) (fileContent *github.RepositoryContent, directoryContent []*github.RepositoryContent, resp *github.Response, err error)
}

type GitHubClientInterface interface {
	GetRepo() GitHubRepositoriesInterface
	SearchCode(repo, query string) ([]string, error)
}

type GitHubClient struct {
	client *github.Client
}

func NewGitHubClient(
	cfg configservice.ConfigServiceInterface,
	kr keyringservice.KeyringServiceInterface,
) GitHubClientInterface {
	serviceName := "gh:github.com"

	token := cfg.GetGithubToken()
	if token == "" {
		var tokenErr error
		token, tokenErr = kr.Get(serviceName, cfg.GetGithubUser())
		if tokenErr != nil {
			panic(tokenErr)
		}
	}

	ctx := context.Background()
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})
	tc := oauth2.NewClient(ctx, ts)

	return &GitHubClient{client: github.NewClient(tc)}
}

func (gh *GitHubClient) GetRepo() GitHubRepositoriesInterface {
	return gh.client.Repositories
}

func (gh *GitHubClient) SearchCode(repo, query string) ([]string, error) {
	q := fmt.Sprintf("repo:%s %s", repo, query)
	codeResult, res, searchErr := gh.client.Search.Code(context.Background(), q, nil)
	if searchErr != nil {
		return nil, searchErr
	}

	if res.StatusCode != 200 {
		body, err := io.ReadAll(res.Body)
		if err != nil {
			return nil, fmt.Errorf("failed to read response body: %w", err)
		}
		defer res.Body.Close()

		return nil, fmt.Errorf("failed to search code: %d. message %s", res.StatusCode, body)
	}

	result := make([]string, len(codeResult.CodeResults))
	for i, code := range codeResult.CodeResults {
		result[i] = code.GetPath()
	}

	return result, nil
}
