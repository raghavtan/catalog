//go:build wireinject

package bind

import (
	"github.com/google/wire"
	"github.com/motain/of-catalog/internal/modules/component/handler"
	"github.com/motain/of-catalog/internal/modules/component/repository"
	"github.com/motain/of-catalog/internal/services/compassservice"
	"github.com/motain/of-catalog/internal/services/configservice"
	"github.com/motain/of-catalog/internal/services/githubservice"
	"github.com/motain/of-catalog/internal/services/keyringservice"
	"github.com/motain/of-catalog/internal/services/prometheusservice"
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

	// Prometheusservice
	prometheusservice.NewPrometheusService,
	prometheusservice.NewPrometheusClient,
	wire.Bind(new(prometheusservice.PrometheusServiceInterface), new(*prometheusservice.PrometheusService)),

	// --- metric module ---
	// Repository
	repository.NewRepository,
	wire.Bind(new(repository.RepositoryInterface), new(*repository.Repository)),

	// BindHandler
	handler.NewBindHandler,
)

func initializeHandler() *handler.BindHandler {
	panic(wire.Build(ProviderSet))
}
