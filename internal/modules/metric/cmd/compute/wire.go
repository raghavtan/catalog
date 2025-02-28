//go:build wireinject

package compute

import (
	"github.com/google/wire"
	"github.com/motain/fact-collector/internal/modules/metric/handler"
	"github.com/motain/fact-collector/internal/modules/metric/handler/factcollectors"
	"github.com/motain/fact-collector/internal/modules/metric/handler/factinterpreter"
	"github.com/motain/fact-collector/internal/modules/metric/repository"
	"github.com/motain/fact-collector/internal/services/compassservice"
	"github.com/motain/fact-collector/internal/services/configservice"
	"github.com/motain/fact-collector/internal/services/githubservice"
	"github.com/motain/fact-collector/internal/services/jsonservice"
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
	githubservice.NewGitHubService,
	wire.Bind(new(githubservice.GitHubServiceInterface), new(*githubservice.GitHubService)),

	// JSONService
	jsonservice.NewJSONService,

	// --- metric module ---
	// Repository
	repository.NewRepository,
	wire.Bind(new(repository.RepositoryInterface), new(*repository.Repository)),

	// FactColletors
	// GithubFactCollector
	factcollectors.NewGithubFactCollector,
	wire.Bind(new(factcollectors.GithubFactCollectorInterface), new(*factcollectors.GithubFactCollector)),

	// JSONFactCollector
	factcollectors.NewJSONAPIFactCollector,
	wire.Bind(new(factcollectors.JSONAPIFactCollectorInterface), new(*factcollectors.JSONAPIFactCollector)),

	// FactInterpreter
	factinterpreter.NewFactInterpreter,
	wire.Bind(new(factinterpreter.FactInterpreterInterface), new(*factinterpreter.FactInterpreter)),

	// ComputeHandler
	handler.NewComputeHandler,
)

func initializeHandler() *handler.ComputeHandler {
	panic(wire.Build(ProviderSet))
}
