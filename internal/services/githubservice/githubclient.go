package githubservice

//go:generate mockgen -destination=./mocks/mock_github_client.go -package=githubservice github.com/motain/of-catalog/internal/services/githubservice GitHubRepositoriesInterface

import (
	"context"
	"fmt"
	"io"
	"time"

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
	return &rateLimitedRepositories{gh.client.Repositories}
}

// rateLimitedRepositories wraps the GitHub repositories client with rate limit handling
type rateLimitedRepositories struct {
	repos *github.RepositoriesService
}

func (r *rateLimitedRepositories) Get(ctx context.Context, owner, repo string) (*github.Repository, *github.Response, error) {
	return executeWithRetry(func() (*github.Repository, *github.Response, error) {
		return r.repos.Get(ctx, owner, repo)
	})
}

func (r *rateLimitedRepositories) GetContents(ctx context.Context, owner, repo, path string, opts *github.RepositoryContentGetOptions) (fileContent *github.RepositoryContent, directoryContent []*github.RepositoryContent, resp *github.Response, err error) {
	return executeWithRetryContents(func() (*github.RepositoryContent, []*github.RepositoryContent, *github.Response, error) {
		return r.repos.GetContents(ctx, owner, repo, path, opts)
	})
}

func (gh *GitHubClient) SearchCode(repo, query string) ([]string, error) {
	q := fmt.Sprintf("repo:%s %s", repo, query)

	codeResult, res, searchErr := executeWithRetrySearch(func() (*github.CodeSearchResult, *github.Response, error) {
		return gh.client.Search.Code(context.Background(), q, nil)
	})

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

// Generic retry function for repository operations
func executeWithRetry[T any](fn func() (*T, *github.Response, error)) (*T, *github.Response, error) {
	maxRetries := 3
	baseDelay := time.Second * 2

	for attempt := 0; attempt < maxRetries; attempt++ {
		result, response, err := fn()

		if err == nil {
			return result, response, nil
		}

		// Check if it's a rate limit error
		if rateLimitErr, ok := err.(*github.RateLimitError); ok {
			if attempt < maxRetries-1 {
				// Wait until rate limit resets, with a small buffer
				waitTime := time.Until(rateLimitErr.Rate.Reset.Time) + time.Second*5
				fmt.Printf("Rate limit hit, waiting %v until reset...\n", waitTime)
				time.Sleep(waitTime)
				continue
			}
		}

		// Check for secondary rate limit (abuse detection)
		if abuseErr, ok := err.(*github.AbuseRateLimitError); ok {
			if attempt < maxRetries-1 {
				waitTime := time.Duration(abuseErr.GetRetryAfter()) * time.Second
				fmt.Printf("Secondary rate limit hit, waiting %v...\n", waitTime)
				time.Sleep(waitTime)
				continue
			}
		}

		// For other errors or final attempt, return the error
		if attempt == maxRetries-1 {
			return result, response, err
		}

		// Exponential backoff for other errors
		delay := baseDelay * time.Duration(1<<attempt)
		fmt.Printf("Request failed (attempt %d/%d), retrying in %v: %v\n", attempt+1, maxRetries, delay, err)
		time.Sleep(delay)
	}

	return nil, nil, fmt.Errorf("max retries exceeded")
}

// Specialized retry function for GetContents (different return signature)
func executeWithRetryContents(fn func() (*github.RepositoryContent, []*github.RepositoryContent, *github.Response, error)) (*github.RepositoryContent, []*github.RepositoryContent, *github.Response, error) {
	maxRetries := 3
	baseDelay := time.Second * 2

	for attempt := 0; attempt < maxRetries; attempt++ {
		fileContent, dirContent, response, err := fn()

		if err == nil {
			return fileContent, dirContent, response, nil
		}

		// Check if it's a rate limit error
		if rateLimitErr, ok := err.(*github.RateLimitError); ok {
			if attempt < maxRetries-1 {
				waitTime := time.Until(rateLimitErr.Rate.Reset.Time) + time.Second*5
				fmt.Printf("Rate limit hit, waiting %v until reset...\n", waitTime)
				time.Sleep(waitTime)
				continue
			}
		}

		// Check for secondary rate limit
		if abuseErr, ok := err.(*github.AbuseRateLimitError); ok {
			if attempt < maxRetries-1 {
				waitTime := abuseErr.GetRetryAfter() * time.Second
				fmt.Printf("Secondary rate limit hit, waiting %v...\n", waitTime)
				time.Sleep(waitTime)
				continue
			}
		}

		// For 404 errors (file not found), don't retry
		if githubErr, ok := err.(*github.ErrorResponse); ok && githubErr.Response.StatusCode == 404 {
			return fileContent, dirContent, response, err
		}

		if attempt == maxRetries-1 {
			return fileContent, dirContent, response, err
		}

		delay := baseDelay * time.Duration(1<<attempt)
		fmt.Printf("Request failed (attempt %d/%d), retrying in %v: %v\n", attempt+1, maxRetries, delay, err)
		time.Sleep(delay)
	}

	return nil, nil, nil, fmt.Errorf("max retries exceeded")
}

// Specialized retry function for Search operations
func executeWithRetrySearch(fn func() (*github.CodeSearchResult, *github.Response, error)) (*github.CodeSearchResult, *github.Response, error) {
	maxRetries := 3
	baseDelay := time.Second * 2

	for attempt := 0; attempt < maxRetries; attempt++ {
		result, response, err := fn()

		if err == nil {
			return result, response, nil
		}

		// Check if it's a rate limit error
		if rateLimitErr, ok := err.(*github.RateLimitError); ok {
			if attempt < maxRetries-1 {
				waitTime := time.Until(rateLimitErr.Rate.Reset.Time) + time.Second*5
				fmt.Printf("Search rate limit hit, waiting %v until reset...\n", waitTime)
				time.Sleep(waitTime)
				continue
			}
		}

		// Check for secondary rate limit
		if abuseErr, ok := err.(*github.AbuseRateLimitError); ok {
			if attempt < maxRetries-1 {
				waitTime := abuseErr.GetRetryAfter() * time.Second
				fmt.Printf("Search secondary rate limit hit, waiting %v...\n", waitTime)
				time.Sleep(waitTime)
				continue
			}
		}

		if attempt == maxRetries-1 {
			return result, response, err
		}

		delay := baseDelay * time.Duration(1<<attempt)
		fmt.Printf("Search request failed (attempt %d/%d), retrying in %v: %v\n", attempt+1, maxRetries, delay, err)
		time.Sleep(delay)
	}

	return nil, nil, fmt.Errorf("max retries exceeded")
}
