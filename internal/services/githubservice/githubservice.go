package githubservice

//go:generate mockgen -destination=./mocks/mock_github_service.go -package=githubservice github.com/motain/of-catalog/internal/services/githubservice GitHubServiceInterface

import (
	"context"
	"fmt"
	"strconv"

	"github.com/google/go-github/v58/github"
)

type GitHubServiceInterface interface {
	GetRepoURL(repo string) string
	GetRepo(repo string) (*github.Repository, error)
	GetFileExists(repo, path string) (bool, error)
	GetFileContent(repo, path string) (string, error)
	GetRepoProperties(repo string) (map[string]string, error)
	Search(repo, query string) ([]string, error)
}

type GitHubService struct {
	client GitHubClientInterface
	owner  string
}

func NewGitHubService(client GitHubClientInterface) *GitHubService {
	return &GitHubService{client: client, owner: "motain"}
}

func (gh *GitHubService) GetRepoURL(repo string) string {
	return fmt.Sprintf("https://github.com/%s/%s", gh.owner, repo)
}

// Get repository details
func (gh *GitHubService) GetRepo(repo string) (*github.Repository, error) {
	ctx := context.Background()
	repository, _, err := gh.client.GetRepo().Get(ctx, gh.owner, repo)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch repo: %w", err)
	}
	return repository, nil
}

func (gh *GitHubService) GetFileExists(repo, path string) (bool, error) {
	ctx := context.Background()
	fileContent, _, _, err := gh.client.GetRepo().GetContents(ctx, gh.owner, repo, path, nil)
	if err != nil {
		if _, ok := err.(*github.ErrorResponse); ok && err.(*github.ErrorResponse).Response.StatusCode == 404 {
			return false, nil
		}
		return false, fmt.Errorf("failed to fetch file: %w", err)
	}

	return fileContent != nil, nil
}

// Get file contents
func (gh *GitHubService) GetFileContent(repo, path string) (string, error) {
	ctx := context.Background()
	fileContent, _, _, fetchErr := gh.client.GetRepo().GetContents(ctx, gh.owner, repo, path, nil)
	if fetchErr != nil {
		return "", fmt.Errorf("failed to fetch file: %w", fetchErr)
	}

	content, decodeErr := fileContent.GetContent()
	if decodeErr != nil {
		return "", fmt.Errorf("failed to decode file content: %w", decodeErr)
	}

	return content, nil
}

func (gh *GitHubService) GetRepoProperties(repo string) (map[string]string, error) {
	ctx := context.Background()

	repoDetails, _, err := gh.client.GetRepo().Get(ctx, gh.owner, repo)
	if err != nil {
		return nil, err
	}

	properties := map[string]string{
		"Name":          repoDetails.GetName(),
		"Description":   repoDetails.GetDescription(),
		"DefaultBranch": repoDetails.GetDefaultBranch(),
		"Visibility":    repoDetails.GetVisibility(),
		"OpenIssues":    strconv.Itoa(repoDetails.GetOpenIssuesCount()),
		"License":       "", // Default empty in case there's no license
	}

	if repoDetails.GetLicense() != nil {
		properties["License"] = repoDetails.GetLicense().GetName()
	}

	return properties, nil
}

func (gh *GitHubService) Search(repo, query string) ([]string, error) {
	repoWithOwner := fmt.Sprintf("%s/%s", gh.owner, repo)
	return gh.client.SearchCode(repoWithOwner, query)
}
