//go:build wireinject

package apply

import (
	"github.com/google/wire"
	"github.com/motain/fact-collector/internal/modules/metric/handler"
	"github.com/motain/fact-collector/internal/modules/metric/repository"
	"github.com/motain/fact-collector/internal/services/compassservice"
	"github.com/motain/fact-collector/internal/services/configservice"
	"github.com/motain/fact-collector/internal/services/githubservice"
	"github.com/motain/fact-collector/internal/services/keyringservice"
)

var ProviderSet = wire.NewSet(
	// Kyeringservice
	keyringservice.NewKeyringService,
	wire.Bind(new(keyringservice.KeyringServiceInterface), new(*keyringservice.KeyringService)),

	// Configservice
	configservice.NewConfigService,
	wire.Bind(new(configservice.ConfigServiceInterface), new(*configservice.ConfigService)),

	// Compassservice
	compassservice.NewGraphQLClient,
	compassservice.NewHTTPClient,
	compassservice.NewCompassService,
	wire.Bind(new(compassservice.CompassServiceInterface), new(*compassservice.CompassService)),

	// Githubservice
	githubservice.NewGitHubClient,
	githubservice.NewGitHubRepositoriesService,
	wire.Bind(new(githubservice.GitHubRepositoriesServiceInterface), new(*githubservice.GitHubRepositoriesService)),

	// --- metric module ---
	// Repository
	repository.NewRepository,
	wire.Bind(new(repository.RepositoryInterface), new(*repository.Repository)),

	// ApplyHandler
	handler.NewApplyHandler,
)

func initializeHandler() *handler.ApplyHandler {
	panic(wire.Build(ProviderSet))
}
