//go:build wireinject

package apply

import (
	"github.com/google/wire"
	"github.com/motain/of-catalog/internal/modules/scorecard/handler"
	"github.com/motain/of-catalog/internal/modules/scorecard/repository"
	"github.com/motain/of-catalog/internal/services/compassservice"
	"github.com/motain/of-catalog/internal/services/configservice"
	"github.com/motain/of-catalog/internal/services/keyringservice"
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
