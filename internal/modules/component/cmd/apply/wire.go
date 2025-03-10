//go:build wireinject

package apply

import (
	"github.com/google/wire"
	"github.com/motain/fact-collector/internal/modules/component/handler"
	"github.com/motain/fact-collector/internal/modules/component/repository"
	"github.com/motain/fact-collector/internal/services/compassservice"
	"github.com/motain/fact-collector/internal/services/configservice"
	"github.com/motain/fact-collector/internal/services/githubservice"
	"github.com/motain/fact-collector/internal/services/keyringservice"
	"github.com/motain/fact-collector/internal/services/ownerservice"
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
	githubservice.NewGitHubService,
	wire.Bind(new(githubservice.GitHubServiceInterface), new(*githubservice.GitHubService)),

	// OwnerService
	ownerservice.NewOwnerService,
	wire.Bind(new(ownerservice.OwnerServiceInterface), new(*ownerservice.OwnerService)),

	// --- component module ---
	// Repository
	repository.NewRepository,
	wire.Bind(new(repository.RepositoryInterface), new(*repository.Repository)),

	// ApplyHandler
	handler.NewApplyHandler,
)

func initializeHandler() *handler.ApplyHandler {
	panic(wire.Build(ProviderSet))
}
