//go:build wireinject

package compute

import (
	"github.com/google/wire"
	"github.com/motain/of-catalog/internal/modules/component/handler"
	"github.com/motain/of-catalog/internal/modules/component/repository"
	"github.com/motain/of-catalog/internal/services/compassservice"
	"github.com/motain/of-catalog/internal/services/configservice"
	"github.com/motain/of-catalog/internal/services/factsystem/aggregators"
	"github.com/motain/of-catalog/internal/services/factsystem/extractors"
	"github.com/motain/of-catalog/internal/services/factsystem/processor"
	"github.com/motain/of-catalog/internal/services/factsystem/validators"
	"github.com/motain/of-catalog/internal/services/githubservice"
	"github.com/motain/of-catalog/internal/services/jsonservice"
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

	// JSONService
	jsonservice.NewJSONService,

	// --- metric module ---
	// Repository
	repository.NewRepository,
	wire.Bind(new(repository.RepositoryInterface), new(*repository.Repository)),
	// Fact System
	aggregators.NewAggregator,
	wire.Bind(new(aggregators.AggregatorInterface), new(*aggregators.Aggregator)),

	extractors.NewExtractor,
	wire.Bind(new(extractors.ExtractorInterface), new(*extractors.Extractor)),

	validators.NewValidator,
	wire.Bind(new(validators.ValidatorInterface), new(*validators.Validator)),

	processor.NewProcessor,
	wire.Bind(new(processor.ProcessorInterface), new(*processor.Processor)),

	// ComputeHandler
	handler.NewComputeHandler,
)

func initializeHandler() *handler.ComputeHandler {
	panic(wire.Build(ProviderSet))
}
