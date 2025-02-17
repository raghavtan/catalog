//go:build wireinject

package app

import (
	"github.com/google/go-github/v58/github"
	"github.com/google/wire"
	"github.com/motain/fact-collector/internal/modules/sample/handler"
	"github.com/motain/fact-collector/internal/modules/sample/repository"
	"github.com/motain/fact-collector/internal/services/configservice"
	"github.com/motain/fact-collector/internal/services/githubservice"
	"github.com/motain/fact-collector/internal/services/keyringservice"
)

var ProviderSet = wire.NewSet(
	keyringservice.NewKeyringService,
	wire.Bind(new(keyringservice.KeyringServiceInterface), new(*keyringservice.KeyringService)),
	configservice.NewConfigService,
	wire.Bind(new(configservice.ConfigServiceInterface), new(*configservice.ConfigService)),
	githubservice.NewGitHubRepositoriesClient,
	wire.Bind(new(githubservice.GitHubRepositoriesInterface), new(*github.RepositoriesService)),
	githubservice.NewGitHubRepositoriesService,
	wire.Bind(new(githubservice.GitHubRepositoriesServiceInterface), new(*githubservice.GitHubRepositoriesService)),
	repository.NewRepository,
	wire.Bind(new(repository.RepositoryInterface), new(*repository.Repository)),
	handler.NewHandler,
)

func InitializeHandler() *handler.Handler {
	panic(wire.Build(ProviderSet))
}
