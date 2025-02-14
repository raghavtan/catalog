package repository

import (
	"fmt"

	"github.com/motain/fact-collector/internal/services/configservice"
	"github.com/motain/fact-collector/internal/services/githubservice"
)

type RepositoryInterface interface {
	FetchData() string
}

type Repository struct {
	config configservice.ConfigServiceInterface
	github githubservice.GitHubServiceInterface
}

func NewRepository(
	cfg configservice.ConfigServiceInterface,
	gh githubservice.GitHubServiceInterface,
) *Repository {
	return &Repository{config: cfg, github: gh}
}

func (r *Repository) FetchData() string {
	repo, repoErr := r.github.GetRepo("motain", "iac")
	if repoErr != nil {
		panic(repoErr)
	}

	fmt.Printf("Repository: %s\n", repo.GetDescription())
	return "Data from GitHub Org: " + r.config.GetGithubOrg()
}
