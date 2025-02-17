package githubservice

import (
	"context"

	"github.com/google/go-github/v58/github"
	"github.com/motain/fact-collector/internal/services/configservice"
	"github.com/motain/fact-collector/internal/services/keyringservice"
	"golang.org/x/oauth2"
)

type GitHubRepositoriesInterface interface {
	Get(ctx context.Context, owner, repo string) (*github.Repository, *github.Response, error)
	GetContents(ctx context.Context, owner, repo, path string, opts *github.RepositoryContentGetOptions) (fileContent *github.RepositoryContent, directoryContent []*github.RepositoryContent, resp *github.Response, err error)
}

func NewGitHubRepositoriesClient(
	cfg configservice.ConfigServiceInterface,
	kr keyringservice.KeyringServiceInterface,
) *github.RepositoriesService {
	serviceName := "gh:github.com"
	token, tokenErr := kr.Get(serviceName, cfg.GetGithubUser())
	if tokenErr != nil {
		panic(tokenErr)
	}

	ctx := context.Background()
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})
	tc := oauth2.NewClient(ctx, ts)

	return github.NewClient(tc).Repositories
}
